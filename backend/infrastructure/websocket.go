package infrastructure

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate specific origins
	},
}

type WebSocketConnection struct {
	Conn *websocket.Conn
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocketConnection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return nil, err
	}
	return &WebSocketConnection{Conn: conn}, nil
}

type ConnectionManager struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	cm := &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
	}
	go cm.startConnectionCleanup(1 * time.Minute)
	return cm
}

func (cm *ConnectionManager) AddConnection(userID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if oldConn, exists := cm.connections[userID]; exists {
		oldConn.Close()
	}
	cm.connections[userID] = conn
}

func (cm *ConnectionManager) RemoveConnection(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if conn, exists := cm.connections[userID]; exists {
		conn.Close()
		delete(cm.connections, userID)
	}
}

func (cm *ConnectionManager) SendIfOnline(userID string, message interface{}) (bool, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, exists := cm.connections[userID]
	if !exists {
		return false, nil
	}

	if err := conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		delete(cm.connections, userID)
		return false, err
	}

	if err := conn.WriteJSON(message); err != nil {
		delete(cm.connections, userID)
		return false, err
	}

	return true, nil
}

func (cm *ConnectionManager) startConnectionCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		cm.mu.Lock()
		for userID, conn := range cm.connections {
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
				conn.Close()
				delete(cm.connections, userID)
				log.Printf("Cleaned up dead connection for %s", userID)
			}
		}
		cm.mu.Unlock()
	}
}