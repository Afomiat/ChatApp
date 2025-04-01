package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
    db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository { // Return pointer
    return &UserRepository{db: db}
}

func (ur *UserRepository) SaveUser(ctx context.Context, user *domain.User) error { // Receiver is pointer
    _, err := ur.db.Collection("users").InsertOne(ctx, user)
    return err
}

func (ur *UserRepository) FindUserByUsername(username string) (*domain.User, error) {
    var user domain.User
    err := ur.db.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }
    return &user, nil
}

func (ur *UserRepository) FindUserByEmail(email string) (*domain.User, error) { 
    var user domain.User
    err := ur.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }
    return &user, nil
}

func (ur *UserRepository) FindAllUsers() ([]domain.User, error) {
    var users []domain.User
    cursor, err := ur.db.Collection("users").Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }

    if err = cursor.All(context.Background(), &users); err != nil {
        return nil, err
    }

    return users, nil
}


func (ur *UserRepository) UpdateUserStatus(userID string, online bool) error {
    filter := bson.M{"_id": userID}
    update := bson.M{"$set": bson.M{"online": online}}

    _, err := ur.db.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// Verify if the user exists after updating the status
	// updateduser, err := ur.FindUserByUsername(userID)
	// if err != nil {
	// 	return fmt.Errorf("failed to verify updated user: %w", err)
	// }
	// fmt.Println("Updated user status.................................: ", updateduser)
    return err
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }

    // Normalize email to lowercase
    email = strings.ToLower(strings.TrimSpace(email))
    filter := bson.M{"email": email}
    
    var user domain.User
    err := ur.db.Collection("users").FindOne(ctx, filter).Decode(&user)
    
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Printf("User not found for email: %s", email)
            return nil, nil
        }
        log.Printf("Database error for email %s: %v", email, err)
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    log.Printf("Found user: %+v", user)
    return &user, nil
}
