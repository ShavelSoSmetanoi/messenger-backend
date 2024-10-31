package models

import "time"

// TODO - возможно придется менять...
type File struct {
	ID         string    `json:"id"`
	MessageID  string    `json:"message_id"`
	SenderID   string    `json:"sender_id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	FileSize   int       `json:"file_size"`
	FilePath   string    `json:"file_path"`
	UploadedAt time.Time `json:"uploaded_at"`
}
