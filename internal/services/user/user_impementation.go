package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Service struct {
	userRepo userDB.UserRepository
}

func NewUserService(repo userDB.UserRepository) *Service {
	return &Service{
		userRepo: repo,
	}
}

// GetSettingsByUserID возвращает настройки пользователя по его ID.
func (h *Service) GetSettingsByUserID(ctx context.Context, userID int) (*models.UserSettings, error) {
	return h.userRepo.GetSettingsByUserID(ctx, userID)
}

// UpdateSettings обновляет тему и цвет сообщений пользователя.
func (h *Service) UpdateSettings(ctx context.Context, userID int, theme, messageColor string) error {
	return h.userRepo.UpdateSettings(ctx, userID, theme, messageColor)
}

// GetUserByID - метод для получения пользователя по ID
func (h *Service) GetUserByID(userID string) (*models.User, error) {
	// Используем репозиторий для получения данных о пользователе
	user, err := h.userRepo.GetUserByID(context.Background(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error fetching user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

// RegisterUser регистрирует нового пользователя
func (h *Service) RegisterUser(c *gin.Context) {
	// Получаем данные из контекста
	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User data not found.",
		})
		return
	}

	// Приводим данные к типу map[string]string
	data := userData.(map[string]string)

	// Теперь у тебя есть данные, например:
	username := data["username"]
	email := data["email"]
	password := data["password"]

	err := h.userRepo.CreateUser(username, email, password, "", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// GetUserProfile возвращает профиль пользователя по его ID
func (h *Service) GetUserProfile(c *gin.Context) {
	// Получаем userID из контекста (например, после авторизации)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Используем репозиторий для получения данных о пользователе
	user, err := h.userRepo.GetUserByID(context.Background(), userID.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Возвращаем данные пользователя в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"uuid":     user.UniqueId,
		"username": user.Username,
		"email":    user.Email,
		"photo":    user.Photo,
		"about":    user.About,
	})
}

// UpdateUserProfile обновляет профиль пользователя
func (h *Service) UpdateUserProfile(c *gin.Context) {
	userID := c.Param("user_id")
	var userUpdate models.UserUpdate

	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userRepo.UpdateUser(context.Background(), userID, userUpdate)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// CheckUserByUsername проверяет наличие пользователя по имени
func (h *Service) CheckUserByUsername(username string) (*models.User, error) {
	user, err := h.userRepo.GetUserByUsername(context.Background(), username)

	if err != nil {
		return nil, err // Возвращаем ошибку в случае проблем с получением пользователя
	}

	return user, nil
}
