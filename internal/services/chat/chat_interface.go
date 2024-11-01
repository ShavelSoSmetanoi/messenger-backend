package chat

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

type ChatServiceInterface interface {
	CreateChat(userID int, name string, participants []string) (*models.Chat, error)
	GetChatsByUserID(userID int) ([]models.Chat, error)
	DeleteChat(chatID string) error
}
