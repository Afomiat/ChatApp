package repository

import (
	"github.com/Afomiat/ChatApp/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
    SaveMessage(message domain.Message) error
    GetMessages(userID string) ([]domain.Message, error)
    FindMessagesBetweenUsers(user1, user2 string) ([]domain.Message, error)
    UpdateUserStatus(userID string, online bool) error
    GetUndeliveredMessages(filter interface{}) (*mongo.Cursor, error)
    FindMessages(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)

}