package Websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

func NotifyUser(userID int, message string) {
	userIDStr := strconv.Itoa(userID)
	conn, ok := Connections[string(userIDStr)]
	if !ok {
		log.Printf("User %d not connected", userID)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error sending message to user %d: %v", userID, err)
		conn.Close()
		delete(Connections, string(userID))
	}
}
