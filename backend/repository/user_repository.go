package repository

import (
	"context"
	"fmt"

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

func (ur *UserRepository) SaveUser(user domain.User) error { // Receiver is pointer
    _, err := ur.db.Collection("users").InsertOne(context.Background(), user)
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