package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/Afomiat/ChatApp/userUtil"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterUsecase struct {
	userRepo *repository.UserRepository // Use pointer
	otpRepo  *repository.OTPRepository
	contextTimeout time.Duration
}

// func (r *RegisterUsecase) GetUserByUserName(username string) (any, any) {
// 	panic("unimplemented")
// }

func NewRegisterUsecase(userRepo *repository.UserRepository, otpRepo *repository.OTPRepository, timeout time.Duration) *RegisterUsecase {
    return &RegisterUsecase{
        userRepo:       userRepo,
        otpRepo:        otpRepo,
        contextTimeout: timeout,
    }
}

func (r *RegisterUsecase) GetUserByUserName(c context.Context , username string) (*domain.User, error) {
	// ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	// defer cancel()
	user, err := r.userRepo.FindUserByUsername(username)

	if err != nil {
		return nil, err

	}

	return user, nil
}

func (r *RegisterUsecase) GetUserByEmail(c context.Context, Email string) (*domain.User, error) {
	// ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	// defer cancel()
	user, err := r.userRepo.FindUserByEmail(Email)

	if err != nil {
		return nil, err

	}

	return user, nil
}
func (r *RegisterUsecase) SendOtp(c context.Context, user *domain.User, smtpusername, smtppassword string) error {
	storedOTP, err := r.otpRepo.GetOtpByEmail(c, user.Email)

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if storedOTP != nil {
		if time.Now().Before(storedOTP.ExpiresAt) {
			return fmt.Errorf("otp already sent")
		}

		if err := r.otpRepo.DeleteOtp(c, storedOTP.Email); err != nil {
			return err
		}

	}

	otp := domain.OTP{
		Value:     userUtil.GenerateOTP(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := r.otpRepo.SaveOTP(c, &otp); err != nil {
		return err
	}

	if err := r.SendEmail(user.Email, otp.Value, smtpusername, smtppassword); err != nil {
		fmt.Print("Error sending email*******************************:", err)
		return err

	}
	return nil
}

func (r *RegisterUsecase) SendEmail(email string, otpValue, smtpusername string, smtppassword string) error {
    from := smtpusername
    password := smtppassword

    to := []string{email}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    message := []byte("Your OTP is " + otpValue)

    auth := smtp.PlainAuth("", from, password, smtpHost)
    val_emal := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	fmt.Println("Email sent:", val_emal)

	return val_emal
}

// In usecase/register_usecase.go
func (r *RegisterUsecase) VerifyOtp(ctx context.Context, otp *domain.VerifyOtp) (*domain.OTP, error) {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
    defer cancel()

    storedOTP, err := r.otpRepo.GetOtpByEmail(ctx, otp.Email)
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    if storedOTP == nil {
        return nil, errors.New("no OTP found for this email")
    }

    if storedOTP.Value != otp.Value {
        return nil, errors.New("invalid OTP")
    }

    if time.Now().After(storedOTP.ExpiresAt) {
        return nil, errors.New("OTP has expired")
    }

    // Delete the OTP after successful verification
    if err := r.otpRepo.DeleteOtp(ctx, storedOTP.Email); err != nil {
        log.Printf("Warning: failed to delete OTP: %v", err)
    }

    return storedOTP, nil
}

func (r *RegisterUsecase) CreateUser(c context.Context, user *domain.User) (*primitive.ObjectID, error) {

	c, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	hashedPassword, err := userUtil.HassPassword(user.Password)

	if err != nil {
		return nil, err
	}

	addUser := domain.User{
		ID: primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
		Online: false,

	}

	err = r.userRepo.SaveUser(c, &addUser)

	return &addUser.ID, err


}