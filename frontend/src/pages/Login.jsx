import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/Login.css';

export default function Login() {
  const [formData, setFormData] = useState({
    email: '',
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
      localStorage.clear();

      const response = await axios.post('http://localhost:8080/login', {
        email: formData.email,
        password: formData.password
      });

      console.log('Login response:', response.data); // Debug log

      // Check for successful response - UPDATED CONDITION
      if (response.status === 200 && response.data.user) {
        // Store user data
        localStorage.setItem('currentUserID', response.data.user.id);
        localStorage.setItem('username', response.data.user.username);
        localStorage.setItem('email', response.data.user.email);
        localStorage.setItem('isAuthenticated', 'true'); // Add this line

        console.log('Navigating to /chat'); // Debug log

        navigate('/chat', { replace: true }); // Added replace option
      } else {
        throw new Error(response.data?.error || 'Login failed');
      }
    } catch (err) {
      console.error('Login error details:', err); // More detailed error log
      localStorage.clear();
      setError(err.response?.data?.error || 'Invalid email or password');
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
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="Email"
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
        <p className="register-link">
          Don't have an account? <span onClick={() => navigate('/register')} className="register-text">Register</span>
        </p>
      </div>
    </div>
  );
}