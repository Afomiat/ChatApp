package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OTPRepository struct {
	db *mongo.Database
}

func NewOtpRepository(db *mongo.Database) *OTPRepository {
	return &OTPRepository{db: db}
}

func (op *OTPRepository) SaveOTP(ctx context.Context,  otp *domain.OTP) error {
	
	_, err := op.db.Collection("otps").InsertOne(ctx, otp)

	return err


}

// In repository/otp_repository.go
func (op *OTPRepository) Initialize(ctx context.Context) error {
    // Create index
    _, err := op.db.Collection("otps").Indexes().CreateOne(
        ctx,
        mongo.IndexModel{
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    )
    return err
}

func (op *OTPRepository) GetOtpByEmail(ctx context.Context, email string) (*domain.OTP, error) {
    // Create context with timeout if none exists
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }

    filter := bson.M{"email": email}
    var otp domain.OTP

    err := op.db.Collection("otps").FindOne(ctx, filter).Decode(&otp)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &otp, nil
}

func (op *OTPRepository) DeleteOtp(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := op.db.Collection("otps").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}