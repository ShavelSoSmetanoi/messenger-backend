package message

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// MessageRepository defines methods for working with messages
type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error)
	IsUserInChat(ctx context.Context, chatID int, userID int) (bool, error)
}

// PostgresMessageRepository is a PostgreSQL implementation of MessageRepository
type PostgresMessageRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresMessageRepository creates a new instance of PostgresMessageRepository
func NewPostgresMessageRepository(db *pgxpool.Pool) *PostgresMessageRepository {
	return &PostgresMessageRepository{DB: db}
}

// CreateMessage creates a new message in the chat
func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	// TODO: Implement the logic for creating a new message
	log.Println("CreateMessage not implemented yet")
	return nil
}

// GetMessagesByChatID retrieves all messages in a chat by chat ID
func (r *PostgresMessageRepository) GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error) {
	// TODO: Implement the logic for retrieving messages by chat ID
	log.Println("GetMessagesByChatID not implemented yet")
	return nil, nil
}

// IsUserInChat checks if a user is a participant in a chat
func (r *PostgresMessageRepository) IsUserInChat(ctx context.Context, chatID int, userID int) (bool, error) {
	// TODO: Implement the logic for checking if a user is in the chat
	log.Println("IsUserInChat not implemented yet")
	return false, nil
}
