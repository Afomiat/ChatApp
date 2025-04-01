package infrastructure

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// In infrastructure/mongo.go
func ConnectMongo(uri string) *mongo.Database {
    clientOptions := options.Client().
        ApplyURI(uri).
        SetConnectTimeout(10 * time.Second).
        SetSocketTimeout(30 * time.Second).
        SetServerSelectionTimeout(10 * time.Second).
        SetMaxPoolSize(50)

    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    log.Println("Successfully connected to MongoDB!")
    return client.Database("chatapp")
}