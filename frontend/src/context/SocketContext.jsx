import { createContext, useState, useEffect, useContext, useRef } from 'react';

const SocketContext = createContext({
  socket: null,
  isConnected: false,
});

export const SocketProvider = ({ children }) => {
  const [socketState, setSocketState] = useState({
    socket: null,
    isConnected: false,
  });
  const wsRef = useRef(null);
  const isMounted = useRef(false);
  const pingIntervalRef = useRef(null);
  const retryCount = useRef(0);

  const setupHeartbeat = (ws) => {
    // Clear any existing interval
    if (pingIntervalRef.current) {
      clearInterval(pingIntervalRef.current);
    }

    // Send ping every 25 seconds (less than server timeout)
    pingIntervalRef.current = setInterval(() => {
      if (ws?.readyState === WebSocket.OPEN) {
        try {
          ws.send(JSON.stringify({ type: 'ping' }));
          console.debug('Sent ping');
        } catch (err) {
          console.error('Ping failed:', err);
        }
      }
    }, 25000); // 25 seconds
  };

  const connectWebSocket = () => {
    const currentUser = localStorage.getItem('currentUserID');
    if (!currentUser) return;

    // Clean up previous connection
    if (wsRef.current) {
      wsRef.current.onopen = null;
      wsRef.current.onclose = null;
      wsRef.current.onerror = null;
      if (wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close(1000, 'Reconnecting');
      }
    }

    try {
      const ws = new WebSocket(`ws://localhost:8080/ws?userID=${currentUser}`);
      wsRef.current = ws;

      ws.onopen = () => {
        if (!isMounted.current) return;
        console.log('WebSocket connected');
        retryCount.current = 0;
        setSocketState({
          socket: ws,
          isConnected: true,
        });
        setupHeartbeat(ws); // Start heartbeat on connection
      };

      ws.onclose = (e) => {
        if (!isMounted.current) return;
        console.log(`WebSocket closed: ${e.code} - ${e.reason}`);
        setSocketState({
          socket: null,
          isConnected: false,
        });

        // Reconnect only for abnormal closures
        if (e.code === 1006) {
          const delay = Math.min(5000, 1000 * Math.pow(2, retryCount.current));
          retryCount.current += 1;
          setTimeout(connectWebSocket, delay);
        }
      };

      ws.onerror = (err) => {
        console.error('WebSocket error:', err);
      };

      // Handle incoming pong messages
      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        if (message.type === 'pong') {
          console.debug('Received pong');
        }
      };

    } catch (err) {
      console.error('WebSocket creation error:', err);
    }
  };

  useEffect(() => {
    isMounted.current = true;
    connectWebSocket();

    return () => {
      isMounted.current = false;
      if (pingIntervalRef.current) {
        clearInterval(pingIntervalRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close(1000, 'Component unmounted');
      }
    };
  }, []);

  return (
    <SocketContext.Provider value={socketState}>
      {children}
    </SocketContext.Provider>
  );
};

export const useSocket = () => useContext(SocketContext);