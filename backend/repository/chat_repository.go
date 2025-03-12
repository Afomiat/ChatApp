package repository

import (
    "context"
    "github.com/Afomiat/ChatApp/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type chatRepository struct {
    collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) ChatRepository {
    return &chatRepository{
        collection: db.Collection("messages"),
    }
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
    _, err := r.collection.InsertOne(ctx, message)
    return err
}

func (r *chatRepository) GetMessages(ctx context.Context) ([]*domain.Message, error) {
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var messages []*domain.Message
    if err = cursor.All(ctx, &messages); err != nil {
        return nil, err
    }
    return messages, nil
}