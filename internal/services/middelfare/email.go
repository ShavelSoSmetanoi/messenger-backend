package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	redis "github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/redis"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// RegisterRequest структура для запроса регистрации
type EmailValidate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// EmailValidator отправляет код на почту и устанавливает таймаут
func EmailValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EmailValidate

		// Извлекаем JSON данные из запроса
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format.",
			})
			c.Abort()
			return
		}

		email := req.Email

		ctx := context.Background()
		if _, err := redis.Rdb.Get(ctx, email).Result(); err == nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Please wait before requesting another code.",
			})
			c.Abort()
			return
		}

		code := "12345" // Генерация случайного кода

		// Логика отправки кода на email (имитация)
		err := sendCodeToEmail(email, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send email.",
			})
			c.Abort()
			return
		}

		// Генерация UUID для пользователя
		uuid := pkg.GenerateUniqueID()

		// Подготовка данных для сохранения в Redis
		userData := map[string]string{
			"username": req.Username,
			"email":    req.Email,
			"password": req.Password,
			"code":     code,
		}

		// Сериализация данных в JSON
		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process registration data.",
			})
			c.Abort()
			return
		}

		// Сохранение данных в Redis с ключом UUID
		err = redis.Rdb.Set(ctx, uuid, userDataJSON, 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save verification code.",
			})
			c.Abort()
			return
		}

		// Возвращаем успешный ответ с UUID
		c.JSON(http.StatusOK, gin.H{
			"UUID": uuid,
		})
		c.Abort()
	}
}

type VerifyCodeRequest struct {
	Code string `json:"code" binding:"required"` // Код верификации, введенный пользователем
	UUID string `json:"uuid" binding:"required"` // UUID пользователя
}

func VerifyCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyCodeRequest

		// Извлекаем JSON данные из запроса
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format.",
			})
			c.Abort()
			return
		}

		ctx := context.Background()

		// Получаем сохраненные данные пользователя из Redis по UUID
		userDataJSON, err := redis.Rdb.Get(ctx, req.UUID).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid UUID or verification code.",
			})
			c.Abort()
			return
		}

		// Распаковка данных из JSON
		var userData map[string]string
		if err := json.Unmarshal([]byte(userDataJSON), &userData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse user data.",
			})
			c.Abort()
			return
		}

		// Проверка кода
		if userData["code"] != req.Code {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid verification code.",
			})
			c.Abort()
			return
		}

		// Код верный, передаем данные пользователя в контекст
		c.Set("userData", userData)

		// Продолжаем выполнение запроса
		c.Next()
	}
}

// Заглушка функции отправки кода на email
func sendCodeToEmail(email, code string) error {
	// Логика отправки email с кодом (например, с использованием внешнего сервиса)
	fmt.Printf("Sent verification code %s to %s\n", code, email)
	return nil
}
