package auth

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService jwtDB.UserTokenRepositoryInterface
}

func NewAuthHandler(repo jwtDB.UserTokenRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		authService: repo,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	//var loginRequest LoginRequest
	//if err := c.ShouldBindJSON(&loginRequest); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

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
	//
	c.JSON(http.StatusOK, gin.H{"token": "idi v pizdy"})
}
