package chat

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatDB"
	"time"
)

type ChatService struct {
	chatRepo chatDB.ChatRepository
}

func NewChatService(repo chatDB.ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: repo,
	}
}

// CreateChat создает новый чат
func (s *ChatService) CreateChat(userID int, name string, participants []string) (*models.Chat, error) {
	// Создаем новый чат
	chat := models.Chat{
		Name:      name,
		CreatedAt: time.Now(),
	}

	return &chat, nil
}

// GetChatsByUserID возвращает все чаты пользователя по его ID
func (s *ChatService) GetChatsByUserID(userID string) ([]models.Chat, error) {
	chats, err := s.chatRepo.GetChatsByUserID(context.Background(), 12)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// DeleteChat удаляет чат по его ID
func (s *ChatService) DeleteChat(chatID string) error {

	return nil
}
