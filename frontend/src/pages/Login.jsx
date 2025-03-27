import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/Login.css';

export default function Login() {
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleLogin = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      // Clear any previous session
      localStorage.removeItem('currentUserID');
      localStorage.removeItem('username');

      const response = await axios.post('http://localhost:8080/login', {
        username: formData.username,
        password: formData.password
      });

      if (!response.data?.user?.id) {
        throw new Error('Invalid server response');
      }

      // Store user data only after successful authentication
      localStorage.setItem('currentUserID', response.data.user.id);
      localStorage.setItem('username', response.data.user.username);

      // Navigate to chat after successful auth
      navigate('/chat');
    } catch (err) {
      localStorage.removeItem('currentUserID');
      localStorage.removeItem('username');
      setError(err.response?.data?.error || 'Login failed. Please try again.');
      console.error('Login error:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Chat App</h1>
        {error && <div className="error-message">{error}</div>}
        <form onSubmit={handleLogin}>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            placeholder="Username"
            required
            disabled={isLoading}
          />
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="Password"
            required
            disabled={isLoading}
          />
          <button 
            type="submit" 
            className="login-button"
            disabled={isLoading}
          >
            {isLoading ? 'Logging in...' : 'Login'}
          </button>
        </form>
      </div>
    </div>
  );
}