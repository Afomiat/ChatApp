package usecase

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ChatUsecase struct {
	Clients map[*websocket.Conn]bool
	Mutex   sync.Mutex
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{
		Clients: make(map[*websocket.Conn]bool),
	}
}

func (cu *ChatUsecase) RegisterClient(conn *websocket.Conn) {
	cu.Mutex.Lock()
	defer cu.Mutex.Unlock()
	cu.Clients[conn] = true
}

func (cu *ChatUsecase) RemoveClient(conn *websocket.Conn) {
	cu.Mutex.Lock()
	defer cu.Mutex.Unlock()
	delete(cu.Clients, conn)
}

func (cu *ChatUsecase) BroadcastMessage(message []byte) {
	cu.Mutex.Lock()
	defer cu.Mutex.Unlock()
	for client := range cu.Clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.Close()
			delete(cu.Clients, client)
		}
	}
}
