import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/register.css';

export default function Register() {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmpass: ''
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

  const handleRegister = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    if (formData.confirmpass !== formData.password) {
      setError('Passwords do not match.');
      setIsLoading(false);
      return;
    }

    try {
      // Clear previous session
      localStorage.clear();

      const response = await axios.post('http://localhost:8080/register', {
        username: formData.username,
        password: formData.password
      });

      // Store only what we need for authentication
      localStorage.setItem('isAuthenticated', 'true');
      localStorage.setItem('username', formData.username);

      // Immediate navigation to chat
      navigate('/chat', { replace: true });

    } catch (err) {
      localStorage.clear();
      setError(err.response?.data?.error || 'Registration failed. Please try again.');
      console.error('Register error:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="register-container">
      <div className="register-card">
        <h1>Chat App</h1>
        {error && <div className="error-message">{error}</div>}
        <form onSubmit={handleRegister}>
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
          <input
            type="password"
            name="confirmpass"
            value={formData.confirmpass}
            onChange={handleChange}
            placeholder="Confirm Password"
            required
            disabled={isLoading}
          />
          <button 
            type="submit" 
            className="register-button"
            disabled={isLoading}
          >
            {isLoading ? 'Registering ...' : 'Register'}
          </button>
        </form>
        <p className="login-link">
          Already have an account? <span onClick={() => navigate('/login')} className="login-text">Login</span>
        </p>
      </div>
    </div>
  );
}