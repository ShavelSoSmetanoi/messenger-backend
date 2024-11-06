package user

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type ServiceInterface interface {
	GetUserProfile(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	CheckUserByUsername(username string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error) // Новый метод для получения пользователя по ID
	RegisterUser(c *gin.Context)
}
