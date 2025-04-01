package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email"  bson:"email" `
	Online   bool               `json:"online" bson:"online"`
}

type AuthLogin struct {
	UserID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username string             `json:"username"`
	Email    string             `json:"email" binding:"required"`
	Password string             `json:"password"`
}