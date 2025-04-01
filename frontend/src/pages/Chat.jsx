import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useSocket } from '../context/SocketContext';
import '../styles/Chat.css';

const Chat = () => {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  const { socket, isConnected } = useSocket();
  const currentUser = localStorage.getItem('currentUserID');
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (!currentUser) {
      navigate('/');
      return;
    }

    const fetchUsers = async () => {
      try {
        setIsLoading(true);
        const response = await axios.get('http://localhost:8080/users');
        setUsers(response.data.users.filter(u => u._id !== currentUser));
      } catch (err) {
        console.error('Failed to fetch users', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchUsers();
  }, [currentUser, navigate]);

  useEffect(() => {
    if (!socket || !selectedUser) return;

    const handleMessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.sender === selectedUser || message.recipient === selectedUser) {
        setMessages(prev => [...prev, message]);
      }
    };

    socket.addEventListener('message', handleMessage);
    return () => {
      if (socket) {
        socket.removeEventListener('message', handleMessage);
      }
    };
  }, [socket, selectedUser]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  useEffect(() => {
    if (!selectedUser) return;

    const fetchConversation = async () => {
      try {
        setIsLoading(true);
        const response = await axios.get(
          `http://localhost:8080/messages?user1=${currentUser}&user2=${selectedUser}&limit=100`
        );
        setMessages(response.data.messages || []);
      } catch (err) {
        console.error('Failed to fetch conversation', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchConversation();
  }, [selectedUser, currentUser]);

  const handleSendMessage = () => {
    if (!isConnected || !newMessage.trim() || !selectedUser) return;

    const message = {
      sender: currentUser,
      recipient: selectedUser,
      content: newMessage,
      timestamp: new Date().toISOString()
    };

    try {
      socket.send(JSON.stringify(message));
      setMessages(prev => [...prev, message]);
      setNewMessage('');
    } catch (err) {
      console.error('Failed to send message:', err);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('currentUserID');
    localStorage.removeItem('username');
    navigate('/');
  };

  return (
    <div className="chat-app">
      <div className="sidebar">
        <div className="user-header">
          <h2>Hello, {localStorage.getItem('username') || 'User'}</h2>
          <div className="connection-status">
            Status: {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
          </div>
          <button onClick={handleLogout} className="logout-button">
            Logout
          </button>
        </div>
        <div className="user-list">
          <h3>Users</h3>
          {isLoading ? (
            <div className="loading">Loading users...</div>
          ) : (
            <ul>
              {users.map(user => (
                <li
                  key={user._id}
                  className={selectedUser === user._id ? 'active' : ''}
                  onClick={() => setSelectedUser(user._id)}
                >
                  <span className={`user-status ${user.online ? 'online' : 'offline'}`}></span>
                  {user.username}
                  {user.online && <span className="online-badge">Online</span>}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      <div className="chat-area">
        {selectedUser ? (
          <>
            <div className="chat-header">
              <h3>Chat with {users.find(u => u._id === selectedUser)?.username || 'User'}</h3>
            </div>
            <div className="messages-container">
              {messages.length === 0 && !isLoading ? (
                <div className="empty-chat">No messages yet. Start the conversation!</div>
              ) : (
                messages.map((msg, i) => (
                  <div
                    key={i}
                    className={`message ${msg.sender === currentUser ? 'sent' : 'received'}`}
                  >
                    <div className="message-content">{msg.content}</div>
                    <div className="message-time">
                      {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </div>
                  </div>
                ))
              )}
              <div ref={messagesEndRef} />
            </div>
            <div className="message-input">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                placeholder={isConnected ? "Type a message..." : "Connecting..."}
                disabled={!isConnected}
              />
              <button 
                onClick={handleSendMessage}
                disabled={!isConnected || !newMessage.trim()}
              >
                Send
              </button>
            </div>
          </>
        ) : (
          <div className="select-user-prompt">
            <p>Select a user to start chatting</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Chat;