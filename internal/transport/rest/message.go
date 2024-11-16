package rest

import (
	"database/sql"
	"encoding/json"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/Websocket"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// InitMessageRouter initialize routes related to message actions
func (h *Handler) InitMessageRouter(r *gin.RouterGroup) {
	r.POST("/chats/:chat_id/messages", h.SendMessageHandler)

	r.GET("/chats/:chat_id/messages", h.GetMessagesHandler)

	r.GET("/chats/:chat_id/messages/last", h.GetLastMessageHandler)

	r.PUT("/chats/:chat_id/messages/:message_id", h.UpdateMessageHandler)

	r.DELETE("/chats/:chat_id/messages/:message_id", h.DeleteMessageHandler)
}

type SendMessageResponse struct {
	Status  string         `json:"status"`
	Message models.Message `json:"message"`
}

type UpdateMessageResponse struct {
	Status    string `json:"status"`
	ChatID    int    `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Content   string `json:"content"`
}

type DeleteMessageResponse struct {
	Status    string `json:"status"`
	ChatID    int    `json:"chat_id"`
	MessageID int    `json:"message_id"`
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
	message, participants, err := h.services.Message.SendMessage(chatIDInt, userID.(string), request.Content, "text")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Отправляем сообщение через WebSocket участникам, кроме отправителя
	for _, participant := range participants {
		if participant.UserID != userID { // Не уведомлять отправителя
			// Сериализуем в JSON

			sendMessage := SendMessageResponse{
				Status:  "send",
				Message: *message,
			}
			jsonData, err := json.Marshal(sendMessage)
			if err != nil {
				log.Printf("Failed to serialize notification: %v", err)
				continue
			}

			// Отправляем JSON уведомление пользователю
			Websocket.NotifyUser(participant.UserID, string(jsonData))
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

// UpdateMessageHandler обновляет содержимое сообщения
func (h *Handler) UpdateMessageHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userIDInt int
	switch id := userID.(type) {
	case int:
		userIDInt = id
	case string:
		var err error
		userIDInt, err = strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	messageIDStr := c.Param("message_id")
	chatIDStr := c.Param("chat_id")

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	participants, err := h.services.Message.UpdateMessage(chatID, userIDInt, messageID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, participant := range participants {
		if participant.UserID != userIDInt {
			response := UpdateMessageResponse{
				Status:    "update",
				ChatID:    chatID,
				MessageID: messageID,
				Content:   req.Content,
			}
			jsonData, err := json.Marshal(response)
			if err != nil {
				log.Printf("Failed to serialize notification: %v", err)
				continue
			}

			Websocket.NotifyUser(participant.UserID, string(jsonData))
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}

// DeleteMessageHandler удаляет сообщение
func (h *Handler) DeleteMessageHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userIDInt int
	switch id := userID.(type) {
	case int:
		userIDInt = id
	case string:
		var err error
		userIDInt, err = strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	messageIDStr := c.Param("message_id")
	chatIDStr := c.Param("chat_id")

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	participants, err := h.services.Message.DeleteMessage(chatID, userIDInt, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, participant := range participants {
		if participant.UserID != userIDInt {
			response := DeleteMessageResponse{
				Status:    "delete",
				ChatID:    chatID,
				MessageID: messageID,
			}
			jsonData, err := json.Marshal(response)
			if err != nil {
				log.Printf("Failed to serialize notification: %v", err)
				continue
			}

			Websocket.NotifyUser(participant.UserID, string(jsonData))
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

func (h *Handler) GetLastMessageHandler(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	message, err := h.services.Message.GetLastMessage(chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No messages found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": message.ID,
		"chat_id":    message.ChatID,
		"user_id":    message.UserID,
		"content":    message.Content,
		"type":       message.Type,
		"created_at": message.CreatedAt,
	})
}
