package models

import "time"

type Chat struct {
	ID int `json:"id"`
	// Tag string : default ""
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
