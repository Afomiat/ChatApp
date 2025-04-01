package domain

import "time"

type OTP struct {
	Value     string    `bson:"value"`
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type VerifyOtp struct {
	Value string `json:"otp" bson:"value"`
	Email string `json:"email" bson:"email"`
}
