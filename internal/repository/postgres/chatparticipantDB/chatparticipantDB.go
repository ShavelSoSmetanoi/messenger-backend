package chatparticipantDB

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChatParticipantRepository defines methods for working with chat participants
type ChatParticipantRepository interface {
	GetChatParticipants(ctx context.Context, chatID int) ([]models.ChatParticipant, error)
}

// PostgresChatParticipantRepository is a PostgreSQL implementation of ChatParticipantRepository
type PostgresChatParticipantRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresChatParticipantRepository creates a new instance of PostgresChatParticipantRepository
func NewPostgresChatParticipantRepository(db *pgxpool.Pool) *PostgresChatParticipantRepository {
	return &PostgresChatParticipantRepository{DB: db}
}

// GetChatParticipants retrieves all participants in a chat by chat ID
func (r *PostgresChatParticipantRepository) GetChatParticipants(ctx context.Context, chatID int) ([]models.ChatParticipant, error) {
	rows, err := r.DB.Query(ctx, "SELECT chat_id, user_id, joined_at FROM chatparticipants WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.ChatParticipant
	for rows.Next() {
		var participant models.ChatParticipant
		if err := rows.Scan(&participant.ChatID, &participant.UserID, &participant.JoinedAt); err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return participants, nil
}
