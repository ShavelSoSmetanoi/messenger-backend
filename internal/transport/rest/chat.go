package rest

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// InitChatRouter initializes routes related to chat operations
func (h *Handler) InitChatRouter(r *gin.RouterGroup) {
	r.POST("/chats", h.CreateChatHandler)                 // Route for creating a chat
	r.GET("/chats", h.GetChatsHandler)                    // Route for fetching chats for a user
	r.GET("/chats/:chat_id/users", h.GetChatUsersHandler) // Route for fetching users of a specific chat
	r.DELETE("/chats/:chat_id", h.DeleteChatHandler)      // Route for deleting a chat
}

// CreateChatRequest defines the structure for the request payload for creating a chat
type CreateChatRequest struct {
	Name         string   `json:"name" binding:"required"`         // Название чата
	Participants []string `json:"participants" binding:"required"` // ID
}

// CreateChatHandler handles the creation of a new chat
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

	chat, err := h.services.Chat.CreateChat(intUserID, request.Name, request.Participants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// GetChatsHandler handles fetching all chats for a specific user
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

	chats, err := h.services.Chat.GetChatsByUserID(intUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// DeleteChatHandler handles deleting a chat
func (h *Handler) DeleteChatHandler(c *gin.Context) {
	chatIDStr := c.Param("chat_id")

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Chat successfully deleted"})
}

// GetChatUsersHandler handles fetching all users in a specific chat
func (h *Handler) GetChatUsersHandler(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat_id"})
		return
	}

	userIDs, err := h.services.Chat.GetChatUserID(chatID)
	if err != nil {
		log.Printf("Error fetching users for chat %d: %v", chatID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_ids": userIDs})
}
