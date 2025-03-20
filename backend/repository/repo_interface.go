package repository

import (
    "context"
    "github.com/Afomiat/ChatApp/domain"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, message *domain.Message) error
	GetMessages(ctx context.Context) ([]*domain.Message, error)
	GetMessagesBetweenUsers(ctx context.Context, user1, user2 string) ([]*domain.Message, error)
	EnsureCollectionExists(ctx context.Context) error
	GetUndeliveredMessages(ctx context.Context, userID string) ([]*domain.Message, error)
	MarkMessageAsDelivered(ctx context.Context, messageID string) error
}