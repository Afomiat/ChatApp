package repository

import (
    "context"
    "github.com/Afomiat/ChatApp/domain"
)

type ChatRepository interface {
    SaveMessage(ctx context.Context, message *domain.Message) error
    GetMessages(ctx context.Context) ([]*domain.Message, error)
}