package routers

import (
	"github.com/Afomiat/ChatApp/delivery/controllers"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database) {
	r.Use(cors.Default())

	chatRepo := repository.NewChatRepository(db)
	userRepo := repository.NewUserRepository(db)

	chatUsecase := usecase.NewChatUsecase(chatRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	chatController := controllers.NewChatController(chatUsecase)
	userController := controllers.NewUserController(userUsecase)

	r.GET("/ws", chatController.HandleWebSocket)

	r.GET("/messages", chatController.GetMessages)

	r.POST("/register", userController.RegisterUser)

	r.GET("/users", userController.GetAllUsers)
	r.POST("/login", userController.LoginUser)

}