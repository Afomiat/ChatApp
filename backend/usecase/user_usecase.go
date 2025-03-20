package usecase

import (
	"context"
	"sync"

	"github.com/Afomiat/ChatApp/domain"
	"golang.org/x/crypto/bcrypt"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/gorilla/websocket"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
	Mutex      sync.Mutex
	Clients    map[*websocket.Conn]string 


}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// RegisterUser registers a new user
func (uc *UserUsecase) RegisterUser(ctx context.Context, user *domain.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Save the user to the database
	return uc.userRepo.RegisterUser(ctx, user)
}

// LoginUser verifies the username and password
func (uc *UserUsecase) LoginUser(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.GetAllUsers(ctx)
}

func (cu *UserUsecase) BroadcastUserList(ctx context.Context) error {
	users, err := cu.userRepo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	cu.Mutex.Lock()
	defer cu.Mutex.Unlock()

	// Broadcast the updated user list to all connected clients
	for conn := range cu.Clients {
		err := conn.WriteJSON(map[string]interface{}{
			"type":     "userList",
			"userList": users,
		})
		if err != nil {
			return err
		}
	}

	return nil
}