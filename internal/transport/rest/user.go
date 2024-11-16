package rest

import (
	"context"
	"database/sql"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) InitUserRouter(r *gin.RouterGroup) {
	r.GET("/profile", h.services.User.GetUserProfile)

	r.PUT("/:user_id", h.services.User.UpdateUserProfile)

	r.GET("/check/:username", h.CheckUserByUsernameHandler)

	r.GET("/user/:user_id", h.GetUserByIDHandler)

	r.GET("/users", h.GetAllUsersHandler)

	r.GET("/settings", h.GetUserSettingsHandler)

	r.PUT("/settings", h.UpdateUserSettingsHandler)
}

func (h *Handler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.services.User.GetAllUsers(c.Request.Context())
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Photo:    user.Photo,
			About:    user.About,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userResponses,
	})
}

// GetUserSettingsHandler возвращает настройки пользователя
func (h *Handler) GetUserSettingsHandler(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Преобразуем userID в int, если это строка
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v // Если userID уже int, просто присваиваем
	case string:
		// Если это строка, пытаемся преобразовать в int
		parsedID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userIDInt = parsedID
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Получаем настройки пользователя из сервиса
	settings, err := h.services.User.GetSettingsByUserID(context.Background(), userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve settings"})
		return
	}

	// Возвращаем настройки пользователя
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateUserSettingsHandler обновляет настройки пользователя
func (h *Handler) UpdateUserSettingsHandler(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Преобразуем userID в int, если это строка
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v // Если userID уже int, просто присваиваем
	case string:
		// Если это строка, пытаемся преобразовать в int
		parsedID, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userIDInt = parsedID
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Структура для получения данных запроса
	var req struct {
		Theme        string `json:"theme"`
		MessageColor string `json:"message_color"`
	}

	// Привязка данных из тела запроса к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Обновление настроек через сервис
	err := h.services.User.UpdateSettings(context.Background(), userIDInt, req.Theme, req.MessageColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

func (h *Handler) CheckUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")

	user, err := h.services.User.CheckUserByUsername(username) // Вызов метода бизнес-логики
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByIDHandler(c *gin.Context) {
	// Получаем user_id из параметров маршрута
	userID := c.Param("user_id")

	// Вызов метода бизнес-логики для получения пользователя по ID
	user, err := h.services.User.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userResponse := models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Photo:    user.Photo,
		About:    user.About,
	}

	c.JSON(http.StatusOK, userResponse)
}
