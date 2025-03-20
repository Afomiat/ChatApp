import React from "react";
import { Routes, Route } from "react-router-dom"; // Remove BrowserRouter import
import Login from "./Login";
import Register from "./Register";
import Chat from "./Chat";

const App = () => {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/chat" element={<Chat />} />
    </Routes>
  );
};

export default App;