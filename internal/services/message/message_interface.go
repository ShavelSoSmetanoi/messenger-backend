package message

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

type ServiceInterface interface {
	SendMessage(chatID int, userID string, content string, typeMsg string) (*models.Message, []models.ChatParticipant, error)
	GetMessages(chatID int, userID int) ([]models.Message, error)
	UpdateMessage(chatID int, userID int, messageID int, content string) ([]models.ChatParticipant, error)
	DeleteMessage(chatID int, userID int, messageID int) ([]models.ChatParticipant, error)
}
