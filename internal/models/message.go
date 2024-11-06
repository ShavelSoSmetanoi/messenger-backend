package models

import "time"

type Message struct {
	ID int `json:"id"`
	// Type [text, file]
	ChatID    int       `json:"chat_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
	ReadAt    time.Time `json:"read_at"`
}
