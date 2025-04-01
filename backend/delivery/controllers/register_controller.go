package controllers

import (
	"fmt"
	"net/http"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/infrastructure"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	registerUsecase *usecase.RegisterUsecase
	env            *infrastructure.Env
}

func NewRegisterController(registerUsecase *usecase.RegisterUsecase) *RegisterController {
	return &RegisterController{
		registerUsecase: registerUsecase,
		env:             infrastructure.NewEnv(),
	}
}

func (uc *RegisterController) RegisterUser(c *gin.Context) {
    var user domain.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    returnUser, _ := uc.registerUsecase.GetUserByUserName(c, user.Username)

    if returnUser != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
        return
    }

    returnUser, _ = uc.registerUsecase.GetUserByEmail(c, user.Email)
    if returnUser != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})

        return
    }

    err := uc.registerUsecase.SendOtp(c, &user, uc.env.SMTPUsername, uc.env.SMTPPassword)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully", "user": user})


  
}

func (uc *RegisterController) Verify(c *gin.Context) {
    var otp domain.VerifyOtp

    if err := c.ShouldBindJSON(&otp); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("Received OTP verification request:", otp)  

    otpResponse, err := uc.registerUsecase.VerifyOtp(c, &otp)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    user := domain.User{
        Username: otpResponse.Username,
        Email:    otpResponse.Email,
        Password: otpResponse.Password,
        Online: false,

    }

    userId, err := uc.registerUsecase.CreateUser(c, &user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "userId": userId})
}