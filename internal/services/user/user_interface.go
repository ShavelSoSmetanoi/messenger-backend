package user

import "github.com/gin-gonic/gin"

type UserServiceInterface interface {
	GetUserProfile(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	CheckUserByUsername(c *gin.Context)
	RegisterUser(c *gin.Context)
}
