package controllers

import (
	"log"

	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"net/http"
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

    // Register the user
    userID := c.Query("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
        return
    }
    cc.Usecase.RegisterClient(conn, userID)
    defer cc.Usecase.RemoveClient(conn)

    // Fetch undelivered messages for the user
    undeliveredMessages, err := cc.Usecase.GetUndeliveredMessages(c.Request.Context(), userID)
    if err != nil {
        log.Println("Failed to fetch undelivered messages:", err)
    } else {
        for _, msg := range undeliveredMessages {
            if err := conn.WriteJSON(map[string]interface{}{
                "type":      "privateMessage",
                "sender":    msg.Sender,
                "content":   msg.Content,
                "timestamp": msg.Timestamp,
            }); err != nil {
                log.Println("Failed to send undelivered message:", err)
                continue
            }

        }
    }

    for {
        var msg struct {
            Type      string `json:"type"`
            SenderID  string `json:"senderId"`
            ReceiverID string `json:"receiverId"`
            Content   string `json:"content"`
        }
        if err := conn.ReadJSON(&msg); err != nil {
            log.Println("WebSocket read error:", err)
            break
        }

        if err := cc.Usecase.SendPrivateMessage(msg.SenderID, msg.ReceiverID, msg.Content); err != nil {
            log.Println("Failed to send private message:", err)
        }
    }
}


func (cc *ChatController) GetMessages(c *gin.Context) {
	user1 := c.Query("user1")
	user2 := c.Query("user2")
	if user1 == "" || user2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user1 and user2 are required"})
		return
	}

	messages, err := cc.Usecase.GetMessagesBetweenUsers(c.Request.Context(), user1, user2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

