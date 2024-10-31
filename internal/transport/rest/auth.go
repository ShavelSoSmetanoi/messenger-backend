package rest

import (
	middleware "github.com/ShavelSoSmetanoi/messenger-backend/internal/services/middelfare"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

func (h *Handler) InitAuthRouter(r *gin.Engine) {
	// Маршруты для аутентификации и авторизации
	r.POST("/verify-email", middleware.EmailValidator())
	r.POST("/register", middleware.VerifyCode(), h.services.User.RegisterUser)

	r.POST("/login", h.HandleLogin)

	// Проверка доступности сервиса
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Вызов бизнес-слоя через `h.services.Auth.Login`
	token, err := h.services.Auth.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Ответ с токеном
	c.JSON(http.StatusOK, gin.H{"token": token})
}
