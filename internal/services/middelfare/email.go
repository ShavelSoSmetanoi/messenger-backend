package middleware

import (
	"context"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

var rdb *redis.Client

// RegisterRequest структура для запроса регистрации
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	About    string `json:"about"`
	Photo    []byte `json:"photo"`
}

// Инициализация Redis
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Адрес Redis
		DB:   0,                // Используемая база
	})

	// Проверяем соединение с Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("could not connect to Redis: %v", err))
	}
}

// EmailValidator отправляет код на почту и устанавливает таймаут
func EmailValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		// Извлекаем JSON данные из запроса
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format.",
			})
			c.Abort()
			return
		}

		email := req.Email // Получаем email из структуры

		// Проверяем, существует ли email в Redis (ожидание повторного запроса)
		ctx := context.Background()
		if _, err := rdb.Get(ctx, email).Result(); err == nil {
			// Если запись существует, отклоняем запрос
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Please wait before requesting another code.",
			})
			c.Abort()
			return
		}

		// Генерация случайного кода
		code := pkg.GenerateCode()

		// Логика отправки кода на email (имитируем отправку)
		err := sendCodeToEmail(email, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send email.",
			})
			c.Abort()
			return
		}

		// Сохранение кода в Redis с таймаутом 5 минут
		err = rdb.Set(ctx, email, code, 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save verification code.",
			})
			c.Abort()
			return
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": "Verification code sent to email.",
		})

		// Прекращаем выполнение (middleware только для отправки кода)
		c.Abort()
	}
}

type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required"` // Email пользователя
	Code  string `json:"code" binding:"required"`  // Код верификации, введенный пользователем
	UUID  string
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
		code, err := rdb.Get(ctx, req.UUID).Result()
		if err != nil || code != req.Code {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid verification code.",
			})
			c.Abort()
			return
		}

		// Код верный, можно продолжать логику (например, активировать пользователя)
		c.JSON(http.StatusOK, gin.H{
			"message": "Code verified successfully.",
		})
	}
}

// Заглушка функции отправки кода на email
func sendCodeToEmail(email, code string) error {
	// Логика отправки email с кодом (например, с использованием внешнего сервиса)
	fmt.Printf("Sent verification code %s to %s\n", code, email)
	return nil
}
