package usecase

import (
	"fmt"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/repository"
)

type UserUsecase struct {
    userRepo *repository.UserRepository // Use pointer
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase { // Accept pointer
    return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) RegisterUser(user domain.User) error {
    return uc.userRepo.SaveUser(user)
}

func (uc *UserUsecase) LoginUser(username, password string) (*domain.User, error) {
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
