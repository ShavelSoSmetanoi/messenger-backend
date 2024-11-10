package messageDB

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// MessageRepository defines methods for working with messages
type MessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error)
	UpdateMessageContent(ctx context.Context, messageID int, content string) error
	DeleteMessage(ctx context.Context, messageID int) error
	IsUserInChat(ctx context.Context, chatID int, userID int) (bool, error)
	IsMessageWrittenByUser(ctx context.Context, messageID int, userID int) (bool, error)
}

// PostgresMessageRepository is a PostgresSQL implementation of MessageRepository
type PostgresMessageRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresMessageRepository creates a new instance of PostgresMessageRepository
func NewPostgresMessageRepository(db *pgxpool.Pool) *PostgresMessageRepository {
	return &PostgresMessageRepository{DB: db}
}

// CreateMessage создает новое сообщение в чате
func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	query := `INSERT INTO messages (chat_id, user_id, type, content, created_at, is_read, read_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.DB.QueryRow(ctx, query,
		message.ChatID,
		message.UserID,
		message.Type,
		message.Content,
		message.CreatedAt,
		message.IsRead,
		message.ReadAt,
	).Scan(&message.ID) // Получение ID нового сообщения

	if err != nil {
		log.Printf("Error creating message: %v", err)
		log.Printf("Message details: ChatID=%d, UserID=%s, Type=%s, Content=%s, CreatedAt=%v, IsRead=%v, ReadAt=%v",
			message.ChatID, message.UserID, message.Type, message.Content, message.CreatedAt, message.IsRead, message.ReadAt)
		return nil, err
	}

	return message, nil
}

// GetMessagesByChatID получает сообщения по ID
func (r *PostgresMessageRepository) GetMessagesByChatID(ctx context.Context, chatID int) ([]models.Message, error) {
	query := `SELECT id, chat_id, user_id, type, content, created_at, is_read, read_at 
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
		if err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.UserID,
			&message.Type, // считываем поле type
			&message.Content,
			&message.CreatedAt,
			&message.IsRead, // считываем поле is_read
			&message.ReadAt, // считываем поле read_at
		); err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// UpdateMessageContent обновляет содержимое сообщения
func (r *PostgresMessageRepository) UpdateMessageContent(ctx context.Context, messageID int, content string) error {
	query := `UPDATE messages SET content = $1 WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, content, messageID)
	if err != nil {
		log.Printf("Error updating message content for MessageID=%d: %v", messageID, err)
		return err
	}
	return nil
}

// DeleteMessage удаляет сообщение из базы данных
func (r *PostgresMessageRepository) DeleteMessage(ctx context.Context, messageID int) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.DB.Exec(ctx, query, messageID)
	if err != nil {
		log.Printf("Error deleting message with MessageID=%d: %v", messageID, err)
		return err
	}
	return nil
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

func (r *PostgresMessageRepository) IsMessageWrittenByUser(ctx context.Context, messageID int, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
                SELECT 1 
                FROM messages 
                WHERE id = $1 AND user_id = $2)`
	err := r.DB.QueryRow(ctx, query, messageID, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if message %d was written by user %d: %v", messageID, userID, err)
		return false, err
	}

	return exists, nil
}
