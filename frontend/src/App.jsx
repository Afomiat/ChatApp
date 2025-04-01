import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Chat from './pages/Chat';
import Register from './pages/Register';
import VerifyOTP from './pages/verifyOTP'
import { SocketProvider } from './context/SocketContext';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/verify" element={<VerifyOTP />} />

        <Route
          path="/chat"
          element={
            <ProtectedRoute>
              <SocketProvider>
                <Chat />
              </SocketProvider>
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

const ProtectedRoute = ({ children }) => {
  const isAuthenticated = localStorage.getItem('isAuthenticated');
  return isAuthenticated ? children : <Navigate to="/" replace />;
};

export default App;