package rest

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/Websocket"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) InitMessageRouter(r *gin.RouterGroup) {

	r.POST("/chats/:chat_id/messages", h.SendMessageHandler)
	r.GET("/chats/:chat_id/messages", h.GetMessagesHandler)
}

func (h *Handler) SendMessageHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chatID := c.Param("chat_id")

	var request struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Используем сервис для отправки сообщения и получения участников
	message, participants, err := h.services.Message.SendMessage(chatIDInt, userID.(string), request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Отправляем сообщение через WebSocket участникам, кроме отправителя
	for _, participant := range participants {
		if participant.UserID != userID { // Не уведомлять отправителя
			Websocket.NotifyUser(participant.UserID, "new_message")
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *Handler) GetMessagesHandler(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Преобразование chatID и userID в int
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string)) // Предполагается, что userID хранится как строка
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// Используем сервис для получения сообщений
	messages, err := h.services.Message.GetMessages(chatID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
