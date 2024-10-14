package middleware

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8" // для Redis
	"golang.org/x/net/context"
	"net/http"
	"time"
)

var rdb = redis.NewClient(&redis.Options{ // Инициализация Redis
	Addr: "localhost:6379", // Адрес Redis
	DB:   0,                // Используемая база
})

// EmailValidator отправляет код на почту и устанавливает таймаут
func EmailValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email") // Получаем email из формы

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
		code := generateCode()

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

// Функция для генерации случайного 6-значного кода
func generateCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%06d", b)
}

// Заглушка функции отправки кода на email
func sendCodeToEmail(email, code string) error {
	// Логика отправки email с кодом (например, с использованием внешнего сервиса)
	fmt.Printf("Sent verification code %s to %s\n", code, email)
	return nil
}
