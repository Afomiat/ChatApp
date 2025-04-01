import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/Login.css';

export default function VerifyOTP() {
  const [otp, setOtp] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [countdown, setCountdown] = useState(60);
  const [resendDisabled, setResendDisabled] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  // Get registration data from navigation state
  const { email, username, password } = location.state || {};

  useEffect(() => {
    if (!email || !username || !password) {
      navigate('/register');
    }

    const timer = countdown > 0 && setInterval(() => {
      setCountdown(countdown - 1);
    }, 1000);

    if (countdown === 0) {
      setResendDisabled(false);
    }

    return () => clearInterval(timer);
  }, [countdown, email, username, password, navigate]);

  const handleVerify = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await axios.post('http://localhost:8080/verify', {
        otp,
        email
      });

      if (response.status === 200) {
        // Auto-login after successful verification
        const loginResponse = await axios.post('http://localhost:8080/login', {
          email,
          password
        });

        if (loginResponse.status === 200 && loginResponse.data.user) {
          localStorage.setItem('currentUserID', loginResponse.data.user.id);
          localStorage.setItem('username', loginResponse.data.user.username);
          localStorage.setItem('email', loginResponse.data.user.email);
          localStorage.setItem('isAuthenticated', 'true');
          navigate('/chat', { replace: true });
        }
      } else {
        throw new Error(response.data?.error || 'Verification failed');
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Invalid OTP. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOTP = async () => {
    try {
      setResendDisabled(true);
      setCountdown(60);
      await axios.post('http://localhost:8080/register', {
        username,
        email,
        password
      });
    } catch (err) {
      setError('Failed to resend OTP. Please try again.');
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Verify Email</h1>
        <p className="otp-instructions">We've sent a verification code to {email}</p>
        {error && <div className="error-message">{error}</div>}
        <form onSubmit={handleVerify}>
          <input
            type="text"
            name="otp"
            value={otp}
            onChange={(e) => setOtp(e.target.value)}
            placeholder="Enter 6-digit OTP"
            required
            disabled={isLoading}
          />
          <button 
            type="submit" 
            className="login-button"
            disabled={isLoading}
          >
            {isLoading ? 'Verifying...' : 'Verify'}
          </button>
        </form>
        <div className="resend-otp">
          {resendDisabled ? (
            <span>Resend OTP in {countdown}s</span>
          ) : (
            <button 
              onClick={handleResendOTP}
              className="resend-button"
            >
              Resend OTP
            </button>
          )}
        </div>
        <p className="register-link">
          Wrong email? <span onClick={() => navigate('/register')} className="register-text">Go back</span>
        </p>
      </div>
    </div>
  );
}