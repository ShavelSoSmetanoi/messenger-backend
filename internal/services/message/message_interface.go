package message

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

type MessageServiceInterface interface {
	SendMessage(chatID int, userID string, content string) (*models.Message, error)
	GetMessages(chatID int, userID int) ([]models.Message, error)
}
