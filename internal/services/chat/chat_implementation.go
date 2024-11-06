package chat

import (
	"context"
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
	chatID, err := s.chatRepo.CreateChat(context.Background(), &chat, participantsIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	// Устанавливаем ID созданного чата
	chat.ID = chatID

	return &chat, nil

}

// GetChatsByUserID возвращает все чаты пользователя по его ID
func (s *ChatService) GetChatsByUserID(userID int) ([]models.Chat, error) {
	// Получаем чаты по userID
	chats, err := s.chatRepo.GetChatsByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// DeleteChat удаляет чат по его ID
func (s *ChatService) DeleteChat(chatID int) error {
	return s.chatRepo.DeleteChat(context.Background(), chatID)
}
