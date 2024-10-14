package user

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/user"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

type UserService struct {
	userRepo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) RegisterUser(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.userRepo.CreateUser(registerRequest.Username, registerRequest.Email, registerRequest.Password, registerRequest.About, registerRequest.Photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

//// GetUserProfile возвращает профиль пользователя по его ID
//func (h *UserService) GetUserProfile(c *gin.Context) {
//	// Получаем userID из контекста (например, после авторизации)
//	userID, exists := c.Get("userID")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
//		return
//	}
//
//	// Используем репозиторий для получения данных о пользователе
//	user, err := h.userRepo.GetUserByID(context.Background(), userID.(string))
//	if err != nil {
//		if err == sql.ErrNoRows {
//			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//			return
//		}
//		log.Printf("Error fetching user by ID: %v", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
//		return
//	}
//
//	// Возвращаем данные пользователя в формате JSON
//	c.JSON(http.StatusOK, gin.H{
//		"id":       user.ID,
//		"username": user.Username,
//		"email":    user.Email,
//		"photo":    user.Photo,
//		"about":    user.About,
//	})
//}
//
//// UpdateUserProfile обновляет профиль пользователя
//func (h *UserService) UpdateUserProfile(c *gin.Context) {
//	userID := c.Param("user_id")
//	var userUpdate models.UserUpdate
//
//	if err := c.ShouldBindJSON(&userUpdate); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	err := h.userRepo.UpdateUser(context.Background(), userID, userUpdate)
//	if err != nil {
//		log.Printf("Error updating user: %v", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
//}
//
//// CheckUserByUsername проверяет наличие пользователя по имени
//func (h *UserService) CheckUserByUsername(c *gin.Context) {
//	username := c.Param("username")
//
//	user, err := h.userRepo.GetUserByUsername(context.Background(), username)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//			return
//		}
//		log.Printf("Error fetching user by username: %v", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
//		return
//	}
//
//	c.JSON(http.StatusOK, user)
//}
