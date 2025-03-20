import React, { useEffect, useState, useContext } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { SocketContext } from "./SocketContext";

const Chat = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedUser, setSelectedUser] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const navigate = useNavigate();
  const socket = useContext(SocketContext); 

  const loggedInUser = JSON.parse(localStorage.getItem("user"));

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await axios.get("http://localhost:8080/users");
        const otherUsers = response.data.filter(
          (user) => user.id !== loggedInUser.id
        );
        setUsers(otherUsers);
        setLoading(false);
      } catch (err) {
        console.error("Failed to fetch users:", err);
        navigate("/");
      }
    };

    fetchUsers();
  }, [navigate, loggedInUser?.id]);

  useEffect(() => {
    if (!socket) return;

    // Listen for incoming messages
    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages((prevMessages) => [...prevMessages, message]);
    };

    // Cleanup listener
    return () => {
      socket.onmessage = null;
    };
  }, [socket]);

  const handleUserClick = async (user) => {
    setSelectedUser(user);
    try {
        const response = await axios.get(
            `http://localhost:8080/messages?user1=${loggedInUser.id}&user2=${user.id}`
        );
        setMessages(response.data || []); 
    } catch (err) {
        console.error("Failed to fetch messages:", err);
        setMessages([]); 
    }
};
  const handleSendMessage = () => {
    if (!newMessage.trim() || !selectedUser || !socket) return;

    const message = {
      type: "privateMessage",
      senderID: loggedInUser.id,
      receiverID: selectedUser.id,
      content: newMessage,
    };

    socket.send(JSON.stringify(message));

    setMessages((prevMessages) => [
      ...prevMessages,
      {
        sender: loggedInUser.id,
        content: newMessage,
        timestamp: new Date().toISOString(),
      },
    ]);
    setNewMessage("");
  };

  if (loading) {
    return <p>Loading users...</p>;
  }

  return (
    <div>
      <h1>Chat</h1>
      <p>Logged in as: {loggedInUser.username}</p>
      <div style={{ display: "flex" }}>
        <div style={{ width: "30%", borderRight: "1px solid #ccc" }}>
          <h2>Users</h2>
          <ul>
            {users.map((user) => (
              <li
                key={user.id}
                onClick={() => handleUserClick(user)}
                style={{ cursor: "pointer" }}
              >
                {user.username} {user.online ? "🟢" : "🔴"}
              </li>
            ))}
          </ul>
        </div>
        <div style={{ width: "70%", padding: "0 10px" }}>
          {selectedUser ? (
            <>
              <h2>Chat with {selectedUser.username}</h2>
              <div
                style={{
                  height: "300px",
                  border: "1px solid #ccc",
                  overflowY: "scroll",
                  padding: "10px",
                }}
              >
                {messages.map((msg, index) => (
                  <div
                    key={index}
                    style={{
                      textAlign:
                        msg.sender === loggedInUser.id ? "right" : "left",
                    }}
                  >
                    <strong>
                      {msg.sender === loggedInUser.id ? "You" : msg.sender}:
                    </strong>{" "}
                    {msg.content}
                    <br />
                    <small>{new Date(msg.timestamp).toLocaleTimeString()}</small>
                  </div>
                ))}
              </div>
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                onKeyPress={(e) => e.key === "Enter" && handleSendMessage()}
                placeholder="Type a message..."
              />
              <button onClick={handleSendMessage}>Send</button>
            </>
          ) : (
            <p>Select a user to start chatting.</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Chat;