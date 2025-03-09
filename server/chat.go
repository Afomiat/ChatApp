package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatRoom struct {
	clients map[*websocket.Conn]string 
	mutex   sync.Mutex
}

var chatRoom = ChatRoom{
	clients: make(map[*websocket.Conn]string),
}

func ChatHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	userID := "User" + fmt.Sprintf("%d", len(chatRoom.clients)+1) // Dummy User IDs for demo
	chatRoom.mutex.Lock()
	chatRoom.clients[conn] = userID
	chatRoom.mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			chatRoom.mutex.Lock()
			delete(chatRoom.clients, conn)
			chatRoom.mutex.Unlock()
			break
		}

		chatRoom.mutex.Lock()
		for client := range chatRoom.clients {
			if client != conn { 
				client.WriteMessage(websocket.TextMessage, []byte(userID+": "+string(msg)))
			}
		}
		chatRoom.mutex.Unlock()
	}
}
