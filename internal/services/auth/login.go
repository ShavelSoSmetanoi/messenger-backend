package auth

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService jwt.UserTokenRepository
}

func NewAuthHandler(authService jwt.UserTokenRepository) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//
	//userID, err := h.authService.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	//	return
	//}
	//
	//token, err := JWT.CreateJWT(userID)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	//	return
	//}

	//c.JSON(http.StatusOK, gin.H{"token": token})
}
