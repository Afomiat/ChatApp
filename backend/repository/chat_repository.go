package repository

// import (
// 	"context"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"github.com/Afomiat/ChatApp/domain"
// )

// type ChatRepository struct {
// 	collection *mongo.Collection
// }

// func NewChatRepository(db *mongo.Database) *ChatRepository {
// 	return &ChatRepository{
// 		collection: db.Collection("messages"),
// 	}
// }

// func (repo *ChatRepository) SaveMessage(ctx context.Context, msg domain.Message) error {
// 	_, err := repo.collection.InsertOne(ctx, msg)
// 	return err
// }
