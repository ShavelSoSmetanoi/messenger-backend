package models

import "time"

// ChatParticipant represents a participant in a chat.
type ChatParticipant struct {
	ChatID   int       `json:"chat_id"`   // ID of the chat to which the user belongs
	UserID   int       `json:"user_id"`   // ID of the user who is a participant in the chat
	JoinedAt time.Time `json:"joined_at"` // The timestamp when the user joined the chat
}
