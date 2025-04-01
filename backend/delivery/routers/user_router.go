package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/Afomiat/ChatApp/delivery/controllers"
    "github.com/Afomiat/ChatApp/repository"
    "github.com/Afomiat/ChatApp/usecase"
    "go.mongodb.org/mongo-driver/mongo"
)

func NewUserRouter(r *gin.Engine, db *mongo.Database) {
    userRepo := repository.NewUserRepository(db) // Returns pointer
    userUsecase := usecase.NewUserUsecase(userRepo) // Pass pointer
    userController := controllers.NewUserController(userUsecase)

    // User registration route
    // r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)
	r.GET("/users", userController.GetAllUsers)
	
}