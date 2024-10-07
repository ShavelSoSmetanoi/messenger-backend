package rest

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// UserTransport - транспортный слой для управления пользователями
type UserTransport struct {
	userService services.UserServiceInterface
}

// NewUserTransport создает новый объект UserTransport
func NewUserTransport(service services.UserServiceInterface) *UserTransport {
	return &UserTransport{
		userService: service,
	}
}

// RegisterRoutes регистрирует маршруты пользователя
func (t *UserTransport) RegisterRoutes(router *gin.Engine) {
	// Маршрут для получения профиля пользователя
	router.GET("/profile", t.GetUserProfile)

	// Маршрут для обновления профиля пользователя
	router.PUT("/profile/:user_id", t.UpdateUserProfile)

	// Маршрут для проверки пользователя по имени пользователя
	router.GET("/check-user/:username", t.CheckUserByUsername)
}

// GetUserProfile - контроллер для получения профиля пользователя
func (t *UserTransport) GetUserProfile(c *gin.Context) {
	t.userService.GetUserProfile(c)
}

// UpdateUserProfile - контроллер для обновления профиля пользователя
func (t *UserTransport) UpdateUserProfile(c *gin.Context) {
	t.userService.UpdateUserProfile(c)
}

// CheckUserByUsername - контроллер для проверки пользователя по имени пользователя
func (t *UserTransport) CheckUserByUsername(c *gin.Context) {
	t.userService.CheckUserByUsername(c)
}
