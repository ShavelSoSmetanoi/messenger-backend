package models

import "time"

// Chat represents a chat entity with its details.
type Chat struct {
	ID        int       `json:"id"`         // Unique identifier for the chat
	Name      string    `json:"name"`       // Name of the chat
	CreatedAt time.Time `json:"created_at"` // The timestamp when the chat was created
}
