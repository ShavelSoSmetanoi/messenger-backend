package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) InitChatRouter(r *gin.RouterGroup) {
	r.POST("/chats", h.CreateChatHandler)
	r.GET("/chats", h.GetChatsHandler)
	// TODO r.DELETE("/chats/:chat_id",...)
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
