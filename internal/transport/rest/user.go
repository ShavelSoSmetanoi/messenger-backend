package rest

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
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
	// Маршрут для регистрации пользователя
	router.POST("/register", t.RegisterUser)

	// Маршрут для получения профиля пользователя
	router.GET("/profile", t.GetUserProfile)

	// Маршрут для обновления профиля пользователя
	router.PUT("/profile/:user_id", t.UpdateUserProfile)

	// Маршрут для проверки пользователя по имени пользователя
	router.GET("/check-user/:username", t.CheckUserByUsername)
}

// RegisterUser - контроллер для регистрации пользователя
func (t *UserTransport) RegisterUser(c *gin.Context) {
	var registerRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		About    string `json:"about"`
		Photo    []byte `json:"photo"`
	}

	// Простой парсинг JSON-запроса без дополнительной логики
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//// Вызов бизнес-слоя для регистрации пользователя
	//if err := t.userService.RegisterUser(c, registerRequest.Username, registerRequest.Email, registerRequest.Password, registerRequest.About, registerRequest.Photo); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// Ответ при успешной регистрации
	c.Status(http.StatusOK) // Возвращает 200 OK без дополнительного тела
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
