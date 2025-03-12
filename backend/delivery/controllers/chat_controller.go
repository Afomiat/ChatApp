package controllers

import (
	"net/http"

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

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			cc.Usecase.RemoveClient(conn)
			break
		}

		cc.Usecase.BroadcastMessage(msg)
	}
}

