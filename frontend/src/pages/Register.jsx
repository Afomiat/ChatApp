import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/Login.css';

export default function Register() {
  const [formData, setFormData] = useState({
    username: '',
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

  const handleRegister = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await axios.post('http://localhost:8080/register', {
        username: formData.username,
        email: formData.email,
        password: formData.password
      });

      if (response.status === 200) {
        navigate('/verify', { 
          state: { 
            email: formData.email,
            username: formData.username,
            password: formData.password
          } 
        });
      } else {
        throw new Error(response.data?.error || 'Registration failed');
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Registration failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Create Account</h1>
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
            {isLoading ? 'Registering...' : 'Register'}
          </button>
        </form>
        <p className="register-link">
          Already have an account? <span onClick={() => navigate('/')} className="register-text">Login</span>
        </p>
      </div>
    </div>
  );
}