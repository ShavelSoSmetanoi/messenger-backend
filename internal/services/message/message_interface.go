package message

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

// ServiceInterface defines the methods for managing messages within a chat application.
type ServiceInterface interface {
	SendMessage(chatID int, userID string, content string, typeMsg string) (*models.Message, []models.ChatParticipant, error)
	GetMessages(chatID int, userID int) ([]models.Message, error)
	UpdateMessage(chatID int, userID int, messageID int, content string) ([]models.ChatParticipant, error)
	DeleteMessage(chatID int, userID int, messageID int) ([]models.ChatParticipant, error)
	GetLastMessage(chatID int) (*models.Message, error)
}
