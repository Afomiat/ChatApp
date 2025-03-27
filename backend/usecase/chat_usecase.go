package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatUsecase struct {
    chatRepo repository.ChatRepository
}

func NewChatUsecase(chatRepo repository.ChatRepository) *ChatUsecase {
    return &ChatUsecase{chatRepo: chatRepo}
}

func (uc *ChatUsecase) HandleMessage(message domain.Message) error {
    return uc.chatRepo.SaveMessage(message)
}

func (uc *ChatUsecase) GetMessagesBetweenUsers(user1, user2 string) ([]domain.Message, error) {
    return uc.chatRepo.FindMessagesBetweenUsers(user1, user2)
}

func (uc *ChatUsecase) UpdateUserStatus(userID string, online bool) error {
    fmt.Println("handler .. enterd *********************")

    return uc.chatRepo.UpdateUserStatus(userID, online)
}

func (uc *ChatUsecase) GetUndeliveredMessages(userID string) ([]domain.Message, error) {
    filter := bson.M{
        "recipient": userID,
        "delivered": false,
        // Optional: Add time window
        "timestamp": bson.M{"$gt": time.Now().Add(-30 * 24 * time.Hour)},
    }

    cursor, err := uc.chatRepo.GetUndeliveredMessages(filter)
    if err != nil {
        return nil, fmt.Errorf("failed to query messages: %w", err)
    }
    defer cursor.Close(context.Background())

    var messages []domain.Message
    if err := cursor.All(context.Background(), &messages); err != nil {
        return nil, fmt.Errorf("failed to decode messages: %w", err)
    }

    return messages, nil
}

func (uc *ChatUsecase) GetConversation(user1, user2 string, limit int) ([]domain.Message, error) {
  
    filter := bson.M{
        "$or": []bson.M{
            {"sender": user1, "recipient": user2},
            {"sender": user2, "recipient": user1},
        },
    }
    
    opts := options.Find().
        SetSort(bson.D{{Key: "timestamp", Value: -1}}). // Newest first
        SetLimit(int64(limit))
    
    cursor, err := uc.chatRepo.FindMessages(filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    var messages []domain.Message
    if err := cursor.All(context.Background(), &messages); err != nil {
        return nil, err
    }

    return messages, nil
}