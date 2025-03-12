package domain

import "time"

type Message struct {
    Sender    string    `json:"sender" bson:"sender"`
    Content   string    `json:"content" bson:"content"`
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}