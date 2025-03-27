package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/Afomiat/ChatApp/domain"
    "github.com/Afomiat/ChatApp/usecase"
    "net/http"
)

type UserController struct {
    userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
    return &UserController{userUsecase: userUsecase}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
    var user domain.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    if err := uc.userUsecase.RegisterUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user})
}

func (uc *UserController) LoginUser(c *gin.Context) {
    var credentials domain.User

    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    user, err := uc.userUsecase.LoginUser(credentials.Username, credentials.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
    users, err := uc.userUsecase.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}