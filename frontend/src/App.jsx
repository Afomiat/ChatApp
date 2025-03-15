import React, { useEffect, useState, useRef } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './App.css';

const App = () => {
  const [messages, setMessages] = useState([]); 
  const [newMessage, setNewMessage] = useState('');
  const ws = useRef(null);
  const username = localStorage.getItem('username'); 
  const navigate = useNavigate();

  useEffect(() => {
    if (!username) {
      navigate('/');
    }
  }, [username, navigate]);

  useEffect(() => {
    axios.get('http://localhost:8080/messages')
      .then(response => {
        setMessages(response.data || []); 
      })
      .catch(error => {
        console.error('Failed to fetch messages:', error);
      });

    ws.current = new WebSocket('ws://localhost:8080/ws');

    ws.current.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.current.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages(prevMessages => [...prevMessages, message]);
    };

    ws.current.onclose = () => {
      console.log('WebSocket connection closed');
    };

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, []);

  const sendMessage = () => {
    if (newMessage.trim() === '') return;

    const message = {
      sender: username, 
      content: newMessage,
      timestamp: new Date().toISOString(),
    };

    ws.current.send(JSON.stringify(message));
    setNewMessage('');
  };

  return (
    <div className="chat-container">
      <div className="messages-container">
        {messages.map((msg, index) => (
          <div
            key={index}
            className={`message ${msg.sender === username ? 'me' : 'other'}`}
          >
            <strong className="strong">{msg.sender}</strong>: {msg.content}
          </div>
        ))}
      </div>
      <div className="input-container">
        <input
          type="text"
          className="input"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
        />
        <button className="button" onClick={sendMessage}>
          Send
        </button>
      </div>
    </div>
  );
};

export default App;