package models

import "time"

// TODO - только для двоих
type Chat struct {
	ID int `json:"id"`
	// Tag string : default ""
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
