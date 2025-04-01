package routers

import (
	"time"

	"github.com/Afomiat/ChatApp/delivery/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Afomiat/ChatApp/repository"
	"github.com/Afomiat/ChatApp/usecase"
)

// In routers/register_router.go
func NewRegisterRouter(r *gin.Engine, db *mongo.Database) {
    userRepo := repository.NewUserRepository(db)
    otpRepo := repository.NewOtpRepository(db)
    
    // Initialize with 5 second timeout
    registerUsecase := usecase.NewRegisterUsecase(userRepo, otpRepo, 5*time.Second)
    
    registerController := controllers.NewRegisterController(registerUsecase)
    r.POST("/register", registerController.RegisterUser)
    r.POST("/verify", registerController.Verify)
}