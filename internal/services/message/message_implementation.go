package message

import (
	"context"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatparticipantDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/messageDB"
	"strconv"
	"time"
)

type Service struct {
	messageRepo         messageDB.MessageRepository
	chatParticipantRepo chatparticipantDB.ChatParticipantRepository
}

func NewMessageService(repo messageDB.MessageRepository, repos chatparticipantDB.ChatParticipantRepository) *Service {
	return &Service{
		messageRepo:         repo,
		chatParticipantRepo: repos,
	}
}

func (h *Service) SendMessage(chatID int, userID string, content string, typeMsg string) (*models.Message, []models.ChatParticipant, error) {
	// Преобразование userID в int
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return nil, nil, err
	}

	// Проверка, является ли пользователь участником чата
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userIDint)
	if err != nil {
		return nil, nil, err
	}
	if !isInChat {
		return nil, nil, fmt.Errorf("доступ предоставлен не будет. свободен")
	}

	message := &models.Message{
		ChatID:    chatID,
		UserID:    userID,
		Content:   content,
		Type:      typeMsg,
		CreatedAt: time.Now(),
	}

	// Сохранение сообщения в базе данных и получение актуальной структуры с ID
	createdMessage, err := h.messageRepo.CreateMessage(context.Background(), message)
	if err != nil {
		return nil, nil, err
	}

	// Получение всех участников чата
	participants, err := h.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, nil, err
	}

	return createdMessage, participants, nil
}

func (h *Service) UpdateMessage(chatID int, userID int, messageID int, content string) ([]models.ChatParticipant, error) {
	// Дополнительная логика (например, проверка прав доступа пользователя на редактирование сообщения)
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("вы не состоите в данном чате")
	}

	// Проверяем, что сообщение написано указанным пользователем
	isAuthor, err := h.messageRepo.IsMessageWrittenByUser(context.Background(), messageID, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке авторства сообщения: %v", err)
	}
	if !isAuthor {
		return nil, fmt.Errorf("у вас нет прав для редактирования этого сообщения")
	}

	// Обновление сообщения в репозитории
	err = h.messageRepo.UpdateMessageContent(context.Background(), messageID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to update message: %v", err)
	}

	// Получение всех участников чата
	participants, err := h.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (h *Service) DeleteMessage(chatID int, userID int, messageID int) ([]models.ChatParticipant, error) {
	// Проверка, является ли пользователь участником чата
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("вы не состоите в данном чате")
	}

	//// Проверка, написал ли пользователь сообщение
	isMessageWrittenByUser, err := h.messageRepo.IsMessageWrittenByUser(context.Background(), messageID, userID)
	if err != nil {
		return nil, err
	}
	if !isMessageWrittenByUser {
		return nil, fmt.Errorf("вы не можете удалить это сообщение, так как не являетесь его автором")
	}

	// Удаление сообщения
	err = h.messageRepo.DeleteMessage(context.Background(), messageID)
	if err != nil {
		return nil, fmt.Errorf("не удалось удалить сообщение: %v", err)
	}

	// Получение всех участников чата
	participants, err := h.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (h *Service) GetMessages(chatID int, userID int) ([]models.Message, error) {
	// Проверка, является ли пользователь участником чата
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("доступ предоставлен не будет. свободен")
	}

	messages, err := h.messageRepo.GetMessagesByChatID(context.Background(), chatID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetLastMessage возврвщает последнее сообщение чата
func (h *Service) GetLastMessage(chatID int) (*models.Message, error) {
	// Вызов репозитория для получения последнего сообщения
	return h.messageRepo.GetLastMessage(context.Background(), chatID)
}
