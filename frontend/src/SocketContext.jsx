import React, { createContext, useEffect, useState } from "react";

const SocketContext = createContext();

const SocketProvider = ({ children }) => {
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    const loggedInUser = JSON.parse(localStorage.getItem("user"));
    if (!loggedInUser) return;

    const newSocket = new WebSocket(`ws://localhost:8080/ws?userID=${loggedInUser.id}`);

    newSocket.onopen = () => {
      console.log("WebSocket connected");
      setSocket(newSocket);
    };

    newSocket.onclose = () => {
      console.log("WebSocket disconnected");
      setSocket(null);
    };

    newSocket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      if (newSocket) {
        newSocket.close();
      }
    };
  }, []);

  return (
    <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
  );
};

export { SocketContext, SocketProvider };