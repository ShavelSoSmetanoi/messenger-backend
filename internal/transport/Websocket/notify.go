package Websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

// NotifyUser sends a message to a specific user over a WebSocket connection.
// It retrieves the connection from the `Connections` map using the user's ID.
// If the user is not connected or an error occurs while sending the message,
// the connection is closed and removed from the map.
func NotifyUser(userID int, message string) {
	userIDStr := strconv.Itoa(userID) // Преобразуем userID в строку
	conn, ok := Connections[userIDStr]
	if !ok {
		log.Printf("User %d not connected", userID)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error sending message to user %d: %v", userID, err)
		conn.Close()
		delete(Connections, userIDStr)
	}
}
