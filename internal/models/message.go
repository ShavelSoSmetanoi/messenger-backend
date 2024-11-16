package models

import "time"

// Message represents a chat message.
type Message struct {
	ID        int       `json:"id"`         // Unique identifier for the message
	Type      string    `json:"type"`       // Type of the message (text, file)
	ChatID    int       `json:"chat_id"`    // ID of the chat the message belongs to
	UserID    string    `json:"user_id"`    // ID of the user who sent the message
	Content   string    `json:"content"`    // The actual content of the message
	CreatedAt time.Time `json:"created_at"` // Timestamp when the message was created
	IsRead    bool      `json:"is_read"`    // Flag indicating whether the message has been read
	ReadAt    time.Time `json:"read_at"`    // Timestamp when the message was read, if applicable
}
