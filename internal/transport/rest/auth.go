package rest

import (
	middleware "github.com/ShavelSoSmetanoi/messenger-backend/internal/services/middelfare"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

// InitAuthRouter sets up routes for authentication and authorization.
func (h *Handler) InitAuthRouter(r *gin.Engine) {
	// Route for email verification.
	r.POST("/verify-email", middleware.EmailValidator())

	// Route for user registration with code verification middleware.
	r.POST("/register", middleware.VerifyCode(), h.services.User.RegisterUser)

	// Route for user login.
	r.POST("/login", h.HandleLogin)

	// Health check endpoint to verify service availability.
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

// LoginRequest represents the request payload for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandleLogin processes user login requests.
func (h *Handler) HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Call the business layer's Login method to authenticate the user.
	token, err := h.services.Auth.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Respond with a success message containing the JWT token.
	c.JSON(http.StatusOK, gin.H{"token": token})
}
