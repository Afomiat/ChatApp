package repository

import (
	"context"
	"log"

	"github.com/Afomiat/ChatApp/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	repo := &UserRepository{
		collection: db.Collection("users"),
	}

	// Ensure the collection exists
	if err := repo.EnsureCollectionExists(context.Background()); err != nil {
		log.Fatalf("Failed to ensure collection exists: %v", err)
	}

	return repo
}

func (r *UserRepository) EnsureCollectionExists(ctx context.Context) error {
	collections, err := r.collection.Database().ListCollectionNames(ctx, bson.M{"name": "users"})
	if err != nil {
		log.Printf("Error listing collections: %v", err)
		return err
	}

	if len(collections) == 0 {
		log.Println("Collection 'users' does not exist. Creating...")
		err := r.collection.Database().CreateCollection(ctx, "users")
		if err != nil {
			log.Printf("Error creating collection: %v", err)
			return err
		}
		log.Println("Collection 'users' created successfully.")
	}

	return nil
}



// RegisterUser stores a new user in the database
func (r *UserRepository) RegisterUser(ctx context.Context, user *domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// GetUserByUsername fetches a user by their username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers fetches all users from the database
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}