package Websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var Connections = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(c *gin.Context) {
	// Извлекаем JWT токен из URL параметра
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Проверка и извлечение userID из токена
	claims, err := parseJWTToken(tokenString, os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID in token"})
		return
	}

	// Обновляем до WebSocket соединения
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Сохраняем соединение для пользователя
	Connections[userID] = conn
	defer delete(Connections, userID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection error: %v", err)
			break
		}
	}
}

func parseJWTToken(tokenString, secret string) (map[string]interface{}, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Проверка валидности токена и извлечение claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
