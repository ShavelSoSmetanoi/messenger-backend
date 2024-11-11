package models

import "time"

type UserSettings struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Theme        string    `json:"theme" db:"theme" default:"light"`
	MessageColor string    `json:"message_color" db:"message_color" default:"#0000FF"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
