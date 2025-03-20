package controllers

import (
	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := primitive.NewObjectID().Hex()

	user := &domain.User{
		ID:       userID,
		Username: req.Username,
		Password: req.Password,
		Online:   false, 
	}

	if err := uc.userUsecase.RegisterUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	if err := uc.userUsecase.BroadcastUserList(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to broadcast user list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": req.Username,
	})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify the username and password
	user, err := uc.userUsecase.LoginUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Return the user ID and username
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
// GetAllUsers fetches all users
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userUsecase.GetAllUsers(c.Request.Context())// what deos the request does >>>>
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

