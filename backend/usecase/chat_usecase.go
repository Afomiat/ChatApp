package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/gorilla/websocket"
)

type ChatUsecase struct {
	Clients    map[*websocket.Conn]string 
	Users      map[string]*websocket.Conn 
	Mutex      sync.Mutex
	Repository repository.ChatRepository

}

func NewChatUsecase(repo repository.ChatRepository) *ChatUsecase {
	return &ChatUsecase{
		Clients:    make(map[*websocket.Conn]string),
		Users:      make(map[string]*websocket.Conn),
		Repository: repo,
	}
}

func (cu *ChatUsecase) RegisterClient(conn *websocket.Conn, userID string) {
	cu.Mutex.Lock()
	defer cu.Mutex.Unlock()
	cu.Clients[conn] = userID
	cu.Users[userID] = conn
}

func (cu *ChatUsecase) RemoveClient(conn *websocket.Conn) {
	cu.Mutex.Lock()
	defer cu.Mutex.Unlock()
	userID := cu.Clients[conn]
	delete(cu.Clients, conn)
	delete(cu.Users, userID)
}



func (cu *ChatUsecase) SendPrivateMessage(senderID, recipientID, message string) error {
    cu.Mutex.Lock()
    defer cu.Mutex.Unlock()

    msg := &domain.Message{
        Sender:    senderID,
        Recipient: recipientID,
        Content:   message,
        Timestamp: time.Now(),
        Delivered: false, 
    }
    if err := cu.Repository.SaveMessage(context.Background(), msg); err != nil {
        return fmt.Errorf("failed to save message: %v", err)
    }

    recipientConn, ok := cu.Users[recipientID]
    if !ok {
        // Recipient is offline, message is saved in the database
        return nil
    }

    // Send the message to the recipient
    if err := recipientConn.WriteJSON(map[string]interface{}{
        "type":    "privateMessage",
        "sender":  senderID,
        "content": message,
    }); err != nil {
        return fmt.Errorf("failed to send message to recipient: %v", err)
    }

    // Mark the message as delivered
    if err := cu.Repository.MarkMessageAsDelivered(context.Background(), msg.ID.Hex()); err != nil {
        return fmt.Errorf("failed to mark message as delivered: %v", err)
    }

    return nil
}
func (cu *ChatUsecase) SaveMessage(ctx context.Context, message *domain.Message) error {
	return cu.Repository.SaveMessage(ctx, message)
}

func (cu *ChatUsecase) GetMessages(ctx context.Context) ([]*domain.Message, error) {
	return cu.Repository.GetMessages(ctx)
}

func (cu *ChatUsecase) GetMessagesBetweenUsers(ctx context.Context, user1, user2 string) ([]*domain.Message, error) {
	return cu.Repository.GetMessagesBetweenUsers(ctx, user1, user2)
}

func (cu *ChatUsecase) GetUndeliveredMessages(ctx context.Context, userID string) ([]*domain.Message, error) {
    return cu.Repository.GetUndeliveredMessages(ctx, userID)
}

// func (cu *ChatUsecase) MarkMessageAsDelivered(ctx context.Context, messageID string) error {
// 	return cu.Repository.MarkMessageAsDelivered(ctx, messageID)
// }