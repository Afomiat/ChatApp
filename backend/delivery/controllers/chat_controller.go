package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/Afomiat/ChatApp/domain"
	"github.com/Afomiat/ChatApp/infrastructure"
	"github.com/Afomiat/ChatApp/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	chatUsecase *usecase.ChatUsecase
	connManager *infrastructure.ConnectionManager
}

func NewChatController(cu *usecase.ChatUsecase, cm *infrastructure.ConnectionManager) *ChatController {
	return &ChatController{
		chatUsecase: cu,
		connManager: cm,
	}
}

func (cc *ChatController) HandleWebSocket(c *gin.Context) {
	ws, err := infrastructure.NewWebSocket(c.Writer, c.Request)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		c.JSON(500, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer ws.Conn.Close()

	userID := c.Query("userID")
	if userID == "" {
		c.JSON(400, gin.H{"error": "userID is required"})
		return
	}

    undeliveredMsgs, err := cc.chatUsecase.GetUndeliveredMessages(userID)
    if err != nil {
        log.Printf("Failed to fetch undelivered messages for %s: %v", userID, err)
    }
	if err := cc.chatUsecase.UpdateUserStatus(userID, true); err != nil {
		log.Printf("Status update failed for %s: %v", userID, err)
		c.JSON(500, gin.H{"error": "Status update failed"})
		return
	}

	cc.connManager.AddConnection(userID, ws.Conn)
	defer func() {
		cc.connManager.RemoveConnection(userID)
		cc.chatUsecase.UpdateUserStatus(userID, false)
	}()

    for _, msg := range undeliveredMsgs {
        if err := ws.Conn.WriteJSON(msg); err != nil {
            log.Printf("Failed to send backlog message to %s: %v", userID, err)
            continue
        }
        // Mark as delivered
        msg.Delivered = true
        if err := cc.chatUsecase.HandleMessage(msg); err != nil {
            log.Printf("Failed to update message status: %v", err)
        }
    }

	// Configure connection
	ws.Conn.SetReadLimit(1024 * 1024) // 1MB
	ws.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.Conn.SetPongHandler(func(string) error {
		ws.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg domain.Message
		if err := ws.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close from %s: %v", userID, err)
			}
			break
		}

		// Validate message
		if msg.Recipient == "" {
			log.Printf("Invalid message from %s: missing recipient", userID)
			continue
		}

		// Set metadata
		msg.Sender = userID
		msg.Timestamp = time.Now()

		// Attempt delivery
		delivered, err := cc.connManager.SendIfOnline(msg.Recipient, msg)
		if err != nil {
			log.Printf("Delivery failed to %s: %v", msg.Recipient, err)
		}
		msg.Delivered = delivered

		// Persist message
		if err := cc.chatUsecase.HandleMessage(msg); err != nil {
			log.Printf("Failed to save message from %s: %v", userID, err)
		}
	}
}
func (cc *ChatController) GetMessagesBetweenUsers(c *gin.Context) {
    // Get query parameters
    user1 := c.Query("user1")
    user2 := c.Query("user2")
    limitStr := c.DefaultQuery("limit", "100")
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid limit parameter"})
        return
    }

    // Fetch conversation
    messages, err := cc.chatUsecase.GetConversation(user1, user2, limit)
    if err != nil {
        log.Printf("Failed to get conversation: %v", err)
        c.JSON(500, gin.H{"error": "failed to retrieve messages"})
        return
    }

    c.JSON(200, gin.H{"messages": messages})
}