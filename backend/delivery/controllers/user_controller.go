package controllers

import (
	"log"
	"net/http"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/infrastructure"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
    userUsecase *usecase.UserUsecase
    env         *infrastructure.Env
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
    env := infrastructure.NewEnv()

    return &UserController{
        userUsecase: userUsecase,
        env: env,
    }
}



func (uc *UserController) LoginUser(c *gin.Context) {
    var credentials domain.User

    if err := c.ShouldBindJSON(&credentials); err != nil {
        log.Printf("Invalid login request: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    authLogin := domain.AuthLogin{
        Email:    credentials.Email,
        Password: credentials.Password,
    }

    user, err := uc.userUsecase.AuthenticateUser(c, &authLogin)
    if err != nil {
        log.Printf("Login failed for %s: %v", authLogin.Email, err)
        // Return generic error message for security
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    log.Printf("User %s logged in successfully", user.Email)
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
    users, err := uc.userUsecase.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}