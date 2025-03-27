package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type chatRepository struct {
    db *mongo.Database
}

func NewChatRepository(db *mongo.Database) ChatRepository {
    return &chatRepository{db: db}
}

func (r *chatRepository) SaveMessage(message domain.Message) error {
    _, err := r.db.Collection("messages").InsertOne(context.Background(), message)
    return err
}

func (r *chatRepository) GetMessages(userID string) ([]domain.Message, error) {
    // Implement logic to fetch messages for a user
    return nil, nil
}

func (r *chatRepository) FindMessagesBetweenUsers(user1, user2 string) ([]domain.Message, error) {
    var messages []domain.Message

    // Query to find messages where:
    // (sender = user1 AND recipient = user2) OR (sender = user2 AND recipient = user1)
    filter := bson.M{
        "$or": []bson.M{
            {"sender": user1, "recipient": user2},
            {"sender": user2, "recipient": user1},
        },
    }

    cursor, err := r.db.Collection("messages").Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }

    if err = cursor.All(context.Background(), &messages); err != nil {
        return nil, err
    }

    return messages, nil
}

func (ur *chatRepository) UpdateUserStatus(userID string, online bool) error {
    // Convert userID string to ObjectID
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        fmt.Printf("Invalid userID: %s, error: %v\n", userID, err)
        return err
    }

    filter := bson.M{"_id": objID} // Use ObjectID in the filter
    update := bson.M{"$set": bson.M{"online": online}}

    result, err := ur.db.Collection("users").UpdateOne(context.Background(), filter, update)
    if err != nil {
        fmt.Printf("Error updating user status: %v\n", err)
        return err
    }
    if result.MatchedCount == 0 {
        fmt.Printf("No user found with userID: %s\n", userID)
        return fmt.Errorf("no user found with userID: %s", userID)
    }
    fmt.Printf("User status updated successfully: %s, online: %v\n", userID, online)
    return nil
}

func (r *chatRepository) GetUndeliveredMessages(filter interface{}) (*mongo.Cursor, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Add sorting by timestamp (oldest first)
    opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
    
    return r.db.Collection("messages").Find(ctx, filter, opts)
}

func (r *chatRepository) FindMessages(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return r.db.Collection("messages").Find(ctx, filter, opts...)
}