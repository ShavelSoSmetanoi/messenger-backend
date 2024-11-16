package models

import "time"

// UserSettings represents the settings of a user, including their preferences for theme and message color.
type UserSettings struct {
	ID           int       `json:"id" db:"id"`                                         // Unique identifier for the user settings entry
	UserID       int       `json:"user_id" db:"user_id"`                               // ID of the user this setting belongs to
	Theme        string    `json:"theme" db:"theme" default:"light"`                   // The user's selected theme, default is "light"
	MessageColor string    `json:"message_color" db:"message_color" default:"#0000FF"` // The user's selected message color, default is blue (#0000FF)
	CreatedAt    time.Time `json:"created_at" db:"created_at"`                         // Timestamp when the settings were created
}
