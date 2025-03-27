import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Chat from './pages/Chat';
import { SocketProvider } from './context/SocketContext';

function App() {
  return (
    <BrowserRouter>
      {/* SocketProvider is now inside ProtectedRoute to ensure it only mounts after auth */}
      <Routes>
        <Route path="/" element={<Login />} />
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
  const currentUser = localStorage.getItem('currentUserID');
  return currentUser ? children : <Navigate to="/" replace />;
};

export default App;