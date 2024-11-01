package messageDB

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

// CreateMessage создает новое сообщение в чате
func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	query := `INSERT INTO messages (chat_id, user_id, content, created_at) 
              VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(ctx, query, message.ChatID, message.UserID, message.Content, message.CreatedAt)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		log.Printf("Message details: ChatID=%d, UserID=%s, Content=%s, CreatedAt=%v",
			message.ChatID, message.UserID, message.Content, message.CreatedAt)
		return err
	}

	return nil
}

// TODO GetMessagesByChatID возвращает все сообщения в чате, 50 сообщений лимит
func (r *PostgresMessageRepository) GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error) {
	query := `SELECT id, chat_id, user_id, content, created_at 
              FROM messages WHERE chat_id = $1`
	rows, err := r.DB.Query(ctx, query, chatID)
	if err != nil {
		log.Printf("Error retrieving messages for ChatID=%d: %v", chatID, err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.CreatedAt); err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// IsUserInChat проверяет, является ли пользователь участником чата
func (r *PostgresMessageRepository) IsUserInChat(ctx context.Context, chatID int, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
                SELECT 1 
                FROM chatparticipants 
                WHERE chat_id = $1 AND user_id = $2)`
	err := r.DB.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user %d is in chat %d: %v", userID, chatID, err)
		return false, err
	}

	return exists, nil
}
