package routers

import (
	"github.com/Afomiat/ChatApp/delivery/controllers"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/Afomiat/ChatApp/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database) {
	chatRepo := repository.NewChatRepository(db)
	chatUsecase := usecase.NewChatUsecase(chatRepo)
	chatController := controllers.NewChatController(chatUsecase)

	r.GET("/ws", chatController.HandleWebSocket)
}
