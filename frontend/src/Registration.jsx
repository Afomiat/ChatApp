import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './App.css'; 

const Registration = () => {
  const [username, setUsername] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    if (username.trim() === '') {
      alert('Please enter a valid username.');
      return;
    }
    localStorage.setItem('username', username); 
    navigate('/chat'); 
  };

  return (
    <div className="registration-container">
      <h1>Register</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Enter your name"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <button type="submit">Register</button>
      </form>
    </div>
  );
};

export default Registration;