package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/Afomiat/ChatApp/delivery/controllers"
    "github.com/Afomiat/ChatApp/repository"
    "github.com/Afomiat/ChatApp/usecase"
    "github.com/Afomiat/ChatApp/infrastructure"
    "go.mongodb.org/mongo-driver/mongo"
)

// SetupChatRoutes initializes the WebSocket route for chat functionality.
func SetupChatRoutes(r *gin.Engine, db *mongo.Database) {
    // Initialize the connection manager to track WebSocket connections.
    connManager := infrastructure.NewConnectionManager()
    chatRepo := repository.NewChatRepository(db)

    chatUsecase := usecase.NewChatUsecase(chatRepo)

    // Initialize the chat controller with the chat use case and connection manager.
    chatController := controllers.NewChatController(chatUsecase, connManager)

    // Define the WebSocket route.
    r.GET("/ws", chatController.HandleWebSocket)
    r.GET("/messages", chatController.GetMessagesBetweenUsers)

}