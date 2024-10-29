package auth

import "github.com/gin-gonic/gin"

// AuthHandlerInterface defines the methods for handling authentication-related requests
type AuthHandlerInterface interface {
	Login(c *gin.Context)
}
