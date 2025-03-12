package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/Afomiat/ChatApp/usecase"
)

type ChatController struct {
	Usecase *usecase.ChatUsecase
}

func NewChatController(uc *usecase.ChatUsecase) *ChatController {
	return &ChatController{Usecase: uc}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (cc *ChatController) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
		return
	}
	defer conn.Close()

	cc.Usecase.RegisterClient(conn)
	defer cc.Usecase.RemoveClient(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var message domain.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue 
		}

		message.Timestamp = time.Now()

		if err := cc.Usecase.SaveMessage(context.Background(), &message); err != nil {
			log.Println("Failed to save message:", err)
			continue 
		}

		cc.Usecase.BroadcastMessage(msg)
	}
}