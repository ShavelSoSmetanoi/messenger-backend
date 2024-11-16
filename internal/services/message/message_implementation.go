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

// NewMessageService creates a new instance of the Service struct
func NewMessageService(repo messageDB.MessageRepository, repos chatparticipantDB.ChatParticipantRepository) *Service {
	return &Service{
		messageRepo:         repo,
		chatParticipantRepo: repos,
	}
}

// SendMessage creates a new message in a chat and retrieves chat participants
func (h *Service) SendMessage(chatID int, userID string, content string, typeMsg string) (*models.Message, []models.ChatParticipant, error) {
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return nil, nil, err
	}

	// Check if the user is a participant of the chat
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userIDint)
	if err != nil {
		return nil, nil, err
	}
	if !isInChat {
		return nil, nil, fmt.Errorf("access denied, you are not a chat participant")
	}

	message := &models.Message{
		ChatID:    chatID,
		UserID:    userID,
		Content:   content,
		Type:      typeMsg,
		CreatedAt: time.Now(),
	}

	createdMessage, err := h.messageRepo.CreateMessage(context.Background(), message)
	if err != nil {
		return nil, nil, err
	}

	participants, err := h.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, nil, err
	}

	return createdMessage, participants, nil
}

// UpdateMessage updates the content of an existing message
func (h *Service) UpdateMessage(chatID int, userID int, messageID int, content string) ([]models.ChatParticipant, error) {
	// Check if the user is a participant of the chat
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("you are not a participant of this chat")
	}

	// Verify that the message was written by the specified user
	isAuthor, err := h.messageRepo.IsMessageWrittenByUser(context.Background(), messageID, userID)
	if err != nil {
		return nil, fmt.Errorf("error verifying message authorship: %v", err)
	}
	if !isAuthor {
		return nil, fmt.Errorf("you do not have permission to edit this message")
	}

	err = h.messageRepo.UpdateMessageContent(context.Background(), messageID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to update message: %v", err)
	}

	participants, err := h.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

// DeleteMessage removes a message from a chat
func (h *Service) DeleteMessage(chatID int, userID int, messageID int) ([]models.ChatParticipant, error) {
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("you are not a participant of this chat")
	}

	isMessageWrittenByUser, err := h.messageRepo.IsMessageWrittenByUser(context.Background(), messageID, userID)
	if err != nil {
		return nil, err
	}
	if !isMessageWrittenByUser {
		return nil, fmt.Errorf("you cannot delete this message as you are not the author")
	}

	err = h.messageRepo.DeleteMessage(context.Background(), messageID)
	if err != nil {
		return nil, fmt.Errorf("не удалось удалить сообщение: %v", err)
	}

	participants, err := h.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (h *Service) GetMessages(chatID int, userID int) ([]models.Message, error) {
	isInChat, err := h.messageRepo.IsUserInChat(context.Background(), chatID, userID)
	if err != nil {
		return nil, err
	}
	if !isInChat {
		return nil, fmt.Errorf("access denied, you are not a chat participant")
	}

	messages, err := h.messageRepo.GetMessagesByChatID(context.Background(), chatID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetLastMessage returns the last message of a chat
func (h *Service) GetLastMessage(chatID int) (*models.Message, error) {
	// Вызов репозитория для получения последнего сообщения
	return h.messageRepo.GetLastMessage(context.Background(), chatID)
}
