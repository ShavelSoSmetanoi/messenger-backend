package user

import "github.com/gin-gonic/gin"

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	About    string `json:"about"`
	Photo    []byte `json:"photo"`
}

type UserServiceInterface interface {
	GetUserProfile(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	CheckUserByUsername(c *gin.Context)
	RegisterUser(c *gin.Context)
}
