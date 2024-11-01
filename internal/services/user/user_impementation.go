package user

import (
	"context"
	"database/sql"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type UserService struct {
	userRepo userDB.UserRepository
}

func NewUserService(repo userDB.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

// RegisterUser регистрирует нового пользователя
func (h *UserService) RegisterUser(c *gin.Context) {
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
func (h *UserService) GetUserProfile(c *gin.Context) {
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
		"uuid":     user.UniqueId,
		"username": user.Username,
		"email":    user.Email,
		"photo":    user.Photo,
		"about":    user.About,
	})
}

// UpdateUserProfile обновляет профиль пользователя
func (h *UserService) UpdateUserProfile(c *gin.Context) {
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
func (h *UserService) CheckUserByUsername(username string) (*models.User, error) {
	user, err := h.userRepo.GetUserByUsername(context.Background(), username)

	if err != nil {
		return nil, err // Возвращаем ошибку в случае проблем с получением пользователя
	}

	return user, nil
}
