package chat

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

type ServiceInterface interface {
	CreateChat(userID int, name string, participants []string) (*models.Chat, error)
	GetChatsByUserID(userID int) ([]models.Chat, error)
	GetChatUserID(chatID int) ([]models.ChatParticipant, error)
	DeleteChat(chatID int) error
}
