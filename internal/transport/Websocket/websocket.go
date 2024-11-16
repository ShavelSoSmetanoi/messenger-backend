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

// Connections maintains active WebSocket connections for users.
// The map key is the user ID, and the value is the WebSocket connection.
var Connections = make(map[string]*websocket.Conn)

// upgrader configures the WebSocket upgrader with relaxed origin checks.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler upgrades the HTTP connection to a WebSocket connection and authenticates the user.
func Handler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

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

// parseJWTToken parses and validates a JWT token and returns the claims if successful.
func parseJWTToken(tokenString, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
