package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatDB"
	"strconv"
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
	// Преобразование слайса string в слайс int
	participantsIDs := make([]int, len(participants))
	for i, p := range participants {
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid participant ID: %v", err)
		}
		participantsIDs[i] = id
	}

	// Создаем новый чат
	chat := models.Chat{
		Name:      name,
		CreatedAt: time.Now(),
	}

	// Вызов репозитория для создания чата и добавления участников
	if err := s.chatRepo.CreateChat(context.Background(), &chat, participantsIDs); err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	return &chat, nil
}

// GetChatsByUserID возвращает все чаты пользователя по его ID
func (s *ChatService) GetChatsByUserID(userID string) ([]models.Chat, error) {
	// Преобразуем userID из строки в int
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Получаем чаты по userID
	chats, err := s.chatRepo.GetChatsByUserID(context.Background(), intUserID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// DeleteChat удаляет чат по его ID
func (s *ChatService) DeleteChat(chatID string) error {

	return nil
}
