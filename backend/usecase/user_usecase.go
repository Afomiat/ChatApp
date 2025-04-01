package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/repository"

	"github.com/Afomiat/ChatApp/userUtil"
)

type UserUsecase struct {
    userRepo *repository.UserRepository // Use pointer
    contextTimeout time.Duration
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase { // Accept pointer
    return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) AuthenticateUser(c context.Context, login *domain.AuthLogin) (*domain.User, error) {
    // Set timeout if not already set
    if _, hasDeadline := c.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        c, cancel = context.WithTimeout(c, 5*time.Second)
        defer cancel()
    }

    log.Printf("Attempting authentication for email: %s", login.Email)
    
    user, err := uc.userRepo.GetUserByEmail(c, login.Email)
    if err != nil {
        log.Printf("Authentication failed for %s: %v", login.Email, err)
        return nil, fmt.Errorf("authentication failed")
    }
    
    if user == nil {
        log.Printf("No user found for email: %s", login.Email)
        return nil, fmt.Errorf("invalid credentials")
    }

    log.Printf("Comparing password for user: %s", user.Email)
    if err := userUtil.ComparePassword(user.Password, login.Password); err != nil {
        log.Printf("Password mismatch for user: %s", user.Email)
        return nil, fmt.Errorf("invalid credentials")
    }

    log.Printf("Authentication successful for user: %s", user.Email)
    return user, nil
}


func (uc *UserUsecase) LoginUser(username, password string, email string) (*domain.User, error) {
    user, err := uc.userRepo.FindUserByUsername(username)
    if err != nil {
        return nil, err
    }

    // In a real application, you should hash the password and compare the hashes.
    if user.Password != password {
        return nil, fmt.Errorf("invalid password")
    }

    return user, nil
}

func (uc *UserUsecase) GetAllUsers() ([]domain.User, error) {
    return uc.userRepo.FindAllUsers()
}
func (uc *UserUsecase) UpdateUserStatus(userID string, online bool) error {
    return uc.userRepo.UpdateUserStatus(userID, online)
}


        
