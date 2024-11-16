package rest

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) InitChatRouter(r *gin.RouterGroup) {
	r.POST("/chats", h.CreateChatHandler)

	r.GET("/chats", h.GetChatsHandler)

	r.GET("/chats/:chat_id/users", h.GetChatUsersHandler)

	r.DELETE("/chats/:chat_id", h.DeleteChatHandler)
}

type CreateChatRequest struct {
	Name         string   `json:"name" binding:"required"`         // Название чата
	Participants []string `json:"participants" binding:"required"` // ID
}

func (h *Handler) CreateChatHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	intUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting user ID"})
		return
	}

	var request CreateChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Вызов сервисного слоя для создания чата
	chat, err := h.services.Chat.CreateChat(intUserID, request.Name, request.Participants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

func (h *Handler) GetChatsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	intUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting user ID"})
		return
	}

	// Вызов сервисного слоя для получения чатов
	chats, err := h.services.Chat.GetChatsByUserID(intUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) DeleteChatHandler(c *gin.Context) {
	// Получаем chat_id из параметров маршрута
	chatIDStr := c.Param("chat_id")

	// Преобразуем chatID в int
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		// Если ошибка преобразования, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Вызов метода бизнес-логики для удаления чата
	err = h.services.Chat.DeleteChat(chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			return
		}
		log.Printf("Error deleting chat by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Успешный ответ при удалении
	c.JSON(http.StatusOK, gin.H{"message": "Chat successfully deleted"})
}

func (h *Handler) GetChatUsersHandler(c *gin.Context) {
	// Преобразуем chat_id из строки в int
	chatID, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat_id"})
		return
	}

	// Получаем список ID пользователей по chatID через сервис
	userIDs, err := h.services.Chat.GetChatUserID(chatID)
	if err != nil {
		log.Printf("Error fetching users for chat %d: %v", chatID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_ids": userIDs})
}
