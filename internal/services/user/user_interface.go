package user

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type ServiceInterface interface {
	GetUserProfile(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	CheckUserByUsername(username string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	RegisterUser(c *gin.Context)
	GetSettingsByUserID(userID int) (*models.UserSettings, error)
	UpdateSettings(userID int, theme, messageColor string) error
}
