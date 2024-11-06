package message

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

type ServiceInterface interface {
	SendMessage(chatID int, userID string, content string) (*models.Message, []models.ChatParticipant, error)
	GetMessages(chatID int, userID int) ([]models.Message, error)
}
