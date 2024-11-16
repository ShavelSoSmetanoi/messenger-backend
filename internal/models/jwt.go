package models

import "time"

// UserToken represents an authentication token for a user.
type UserToken struct {
	ID        int       `json:"id"`         // Unique identifier for the token entry
	UserID    string    `json:"user_id"`    // ID of the user to whom the token belongs
	Token     string    `json:"token"`      // The actual authentication token
	CreatedAt time.Time `json:"created_at"` // Timestamp when the token was created
}
