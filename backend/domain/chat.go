package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"` 

	Sender    string    `json:"sender" bson:"sender"`
	Recipient string    `json:"recipient" bson:"recipient"`
	Content   string    `json:"content" bson:"content"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Delivered bool      `json:"delivered" bson:"delivered"` 

}