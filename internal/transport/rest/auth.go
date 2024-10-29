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

	r.POST("/login", h.services.Auth.Login)

	// Проверка доступности сервиса
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}
