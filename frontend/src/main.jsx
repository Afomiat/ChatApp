import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './index.css';
import Registration from './Registration';
import App from './App';

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <Router>
      <Routes>
        <Route path="/" element={<Registration />} /> 
        <Route path="/chat" element={<App />} /> 
      </Routes>
    </Router>
  </StrictMode>
);