package message

import (
	"context"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/messageDB"
	"time"
)

type MessageService struct {
	messageRepo messageDB.MessageRepository
}

func NewMessageService(repo messageDB.MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: repo,
	}
}

func (h *MessageService) SendMessage(chatID int, userID string, content string) (*models.Message, error) {
	message := &models.Message{
		ChatID:    chatID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := h.messageRepo.CreateMessage(context.Background(), message); err != nil {
		return nil, err
	}

	return message, nil
}

func (h *MessageService) GetMessages(chatID int, userID int) ([]models.Message, error) {
	// Проверка, является ли пользователь участником чата
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("access denied")
	}

	messages, err := h.messageRepo.GetMessagesByChatID(context.Background(), chatID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
