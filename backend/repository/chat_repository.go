package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Afomiat/ChatApp/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type chatRepository struct {
    collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) ChatRepository {
    if db == nil {
        log.Fatalf("Error: db is nil")
    }

    collection := db.Collection("messages")
    if collection == nil {
        log.Fatalf("Error: collection is nil")
    }

    repo := &chatRepository{
        collection: collection,
    }

    // Ensure the collection exists
    if err := repo.EnsureCollectionExists(context.Background()); err != nil {
        log.Fatalf("Failed to ensure collection exists: %v", err)
    }

    log.Println("Successfully initialized ChatRepository with messages collection")
    return repo
}

func (r *chatRepository) EnsureCollectionExists(ctx context.Context) error {
    if r.collection == nil {
        return fmt.Errorf("collection is nil")
    }

    collections, err := r.collection.Database().ListCollectionNames(ctx, bson.M{"name": "messages"})
    if err != nil {
        log.Printf("Error listing collections: %v", err)
        return err
    }

    // If the "messages" collection doesn't exist, create it
    if len(collections) == 0 {
        log.Println("Collection 'messages' does not exist. Creating...")
        err := r.collection.Database().CreateCollection(ctx, "messages")
        if err != nil {
            log.Printf("Error creating collection: %v", err)
            return err
        }
        log.Println("Collection 'messages' created successfully.")
    }

    return nil
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
    if r.collection == nil {
        log.Println("Database collection is not initialized.")
        return fmt.Errorf("database collection is not initialized")
    }

    // message.Delivered = false
    result, err := r.collection.InsertOne(ctx, message)

    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        message.ID = oid
    } else {
        return fmt.Errorf("failed to get inserted message ID")
    }
    if err != nil {
        log.Printf("Error saving message to database: %v", err)
        return err
    }

    log.Println("Message saved successfully.")
    return nil
}

func (r *chatRepository) GetMessages(ctx context.Context) ([]*domain.Message, error) {
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        log.Printf("Error fetching messages from database: %v", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    var messages []*domain.Message
    if err = cursor.All(ctx, &messages); err != nil {
        log.Printf("Error decoding messages: %v", err)
        return nil, err
    }

    if messages == nil {
        log.Println("No messages found in the database.")
        return []*domain.Message{}, nil
    }

    log.Printf("Fetched %d messages from the database.", len(messages))
    return messages, nil
}

func (r *chatRepository) GetMessagesBetweenUsers(ctx context.Context, user1, user2 string) ([]*domain.Message, error) {
    filter := bson.M{
        "$or": []bson.M{
            {"sender": user1, "recipient": user2},
            {"sender": user2, "recipient": user1},
        },
    }
    cursor, err := r.collection.Find(ctx, filter)
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

func (r *chatRepository) GetUndeliveredMessages(ctx context.Context, userID string) ([]*domain.Message, error) {
    if r.collection == nil {
        log.Println("Error: collection is nil")
        return nil, fmt.Errorf("collection is nil")
    }

    filter := bson.M{
        "recipient": userID,
        "delivered": false, // Fetch messages where delivered is false
    }
    cursor, err := r.collection.Find(ctx, filter)
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

func (r *chatRepository) MarkMessageAsDelivered(ctx context.Context, messageID string) error {
    // Convert the messageID string to an ObjectID
    oid, err := primitive.ObjectIDFromHex(messageID)
    if err != nil {
        return fmt.Errorf("invalid message ID: %v", err)
    }

    filter := bson.M{"_id": oid}
    update := bson.M{"$set": bson.M{"delivered": true}}
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to mark message as delivered: %v", err)
    }
    if result.ModifiedCount == 0 {
        return fmt.Errorf("no message found with ID %s", messageID)
    }
    return nil
}