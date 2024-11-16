package chat

import (
	"context"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatparticipantDB"
	"strconv"
	"time"
)

// Service struct provides chat-related operations by interacting with the repositories.
type Service struct {
	chatRepo            chatDB.ChatRepository
	chatParticipantRepo chatparticipantDB.ChatParticipantRepository
}

// NewChatService initializes and returns a new instance of the Service struct.
func NewChatService(repo chatDB.ChatRepository, repos chatparticipantDB.ChatParticipantRepository) *Service {
	return &Service{
		chatRepo:            repo,
		chatParticipantRepo: repos,
	}
}

// GetChatUserID retrieves the participants of a chat by its ID.
func (s *Service) GetChatUserID(chatID int) ([]models.ChatParticipant, error) {
	participants, err := s.chatParticipantRepo.GetChatParticipants(context.Background(), chatID)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

// CreateChat creates a new chat and associates participants with it.
func (s *Service) CreateChat(userID int, name string, participants []string) (*models.Chat, error) {
	participantsIDs := make([]int, len(participants))
	for i, p := range participants {
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid participant ID: %v", err)
		}
		participantsIDs[i] = id
	}

	chat := models.Chat{
		Name:      name,
		CreatedAt: time.Now(),
	}

	chatID, err := s.chatRepo.CreateChat(context.Background(), &chat, participantsIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	chat.ID = chatID

	return &chat, nil
}

// GetChatsByUserID retrieves all chats associated with a specific user ID.
func (s *Service) GetChatsByUserID(userID int) ([]models.Chat, error) {
	chats, err := s.chatRepo.GetChatsByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// DeleteChat removes a chat by its ID from the repository.
func (s *Service) DeleteChat(chatID int) error {
	return s.chatRepo.DeleteChat(context.Background(), chatID)
}
