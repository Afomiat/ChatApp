package routers

import (
	"github.com/Afomiat/ChatApp/delivery/controllers"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database) {
	chatUsecase := usecase.NewChatUsecase()
	chatController := controllers.NewChatController(chatUsecase)

	r.GET("/ws", chatController.HandleWebSocket)
}
