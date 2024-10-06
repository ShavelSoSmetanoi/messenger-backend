package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default() // Создаём новый роутер

	// Определяем маршрут
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Запускаем сервер на порту 8080
	r.Run(":8080")
}
