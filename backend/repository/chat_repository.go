package repository

import (
    "context"
    "fmt"
    "log"
    "github.com/Afomiat/ChatApp/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type chatRepository struct {
    collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) ChatRepository {
    repo := &chatRepository{
        collection: db.Collection("messages"),
    }

    // Ensure the collection exists
    if err := repo.EnsureCollectionExists(context.Background()); err != nil {
        log.Fatalf("Failed to ensure collection exists: %v", err)
    }

    return repo
}

func (r *chatRepository) EnsureCollectionExists(ctx context.Context) error {
    collections, err := r.collection.Database().ListCollectionNames(ctx, bson.M{"name": "messages"})
    if err != nil {
        log.Printf("Error listing collections: %v", err)
        return err
    }

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

    _, err := r.collection.InsertOne(ctx, message)
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