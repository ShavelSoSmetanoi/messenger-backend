package models

import "time"

// ChatParticipant represents a participant in a chat
type ChatParticipant struct {
	ChatID   int       `json:"chat_id"`
	UserID   int       `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}
