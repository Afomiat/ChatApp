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
	// Enable CORS
	r.Use(cors.Default())

	// Initialize repository, usecase, and controller
	chatRepo := repository.NewChatRepository(db)
	chatUsecase := usecase.NewChatUsecase(chatRepo)
	chatController := controllers.NewChatController(chatUsecase)

	// WebSocket route for real-time messaging
	r.GET("/ws", chatController.HandleWebSocket)

	// HTTP GET route to fetch all messages
	r.GET("/messages", func(c *gin.Context) {
		messages, err := chatUsecase.GetMessages(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch messages"})
			return
		}
		c.JSON(200, messages)
	})
}
