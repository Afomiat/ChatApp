package usecase

import (
	"context"
	"sync"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/gorilla/websocket"
)

type ChatUsecase struct {
	Clients    map[*websocket.Conn]bool
	Mutex      sync.Mutex
	Repository repository.ChatRepository
}

func NewChatUsecase(repo repository.ChatRepository) *ChatUsecase {
	return &ChatUsecase{
		Clients:    make(map[*websocket.Conn]bool),
		Repository: repo,
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

func (cu *ChatUsecase) SaveMessage(ctx context.Context, message *domain.Message) error {
	return cu.Repository.SaveMessage(ctx, message)
}

func (cu *ChatUsecase) GetMessages(ctx context.Context) ([]*domain.Message, error) {
	return cu.Repository.GetMessages(ctx)
}