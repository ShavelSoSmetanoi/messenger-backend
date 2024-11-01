package user

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type UserServiceInterface interface {
	GetUserProfile(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	CheckUserByUsername(username string) (*models.User, error)
	RegisterUser(c *gin.Context)
}
