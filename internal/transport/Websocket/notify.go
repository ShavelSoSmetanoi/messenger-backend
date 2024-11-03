package Websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

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
		delete(Connections, userIDStr) // Удаляем соединение, используя строковый userID
	}
}
