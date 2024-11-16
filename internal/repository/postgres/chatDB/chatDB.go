package chatDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
	"time"
)

// ChatRepository defines the interface for interacting with the chat-related database operations.
type ChatRepository interface {
	CreateChat(ctx context.Context, chat *models.Chat, participants []int) (int, error)
	GetUserIDsByNicknames(ctx context.Context, nicknames []string) ([]int, error)
	GetChatsByUserID(ctx context.Context, userID int) ([]models.Chat, error)
	DeleteChat(ctx context.Context, chatID int) error
}

// PostgresChatRepository is the implementation of the ChatRepository interface for PostgreSQL.
// It provides methods for chat operations, such as creating, retrieving, and deleting chats, using a PostgreSQL connection pool.
type PostgresChatRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresChatRepository creates a new instance of PostgresChatRepository with the provided PostgreSQL connection pool.
// It initializes the repository to interact with the database using the given connection pool.
func NewPostgresChatRepository(db *pgxpool.Pool) *PostgresChatRepository {
	return &PostgresChatRepository{DB: db}
}

// GetUserIDsByNicknames retrieves the user IDs corresponding to the provided nicknames.
// It generates a dynamic SQL query to select user IDs from the 'users' table based on the given nicknames.
func (r *PostgresChatRepository) GetUserIDsByNicknames(ctx context.Context, nicknames []string) ([]int, error) {
	if len(nicknames) == 0 {
		log.Println("No nicknames provided")
		return nil, errors.New("no nicknames provided")
	}

	// Создаем плейсхолдеры для никнеймов ($1, $2, и т.д.)
	queryPlaceholders := make([]string, len(nicknames))
	args := make([]interface{}, len(nicknames))
	for i, nickname := range nicknames {
		queryPlaceholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = nickname
	}

	query := fmt.Sprintf("SELECT id FROM users WHERE username IN (%s)", strings.Join(queryPlaceholders, ","))
	log.Printf("Generated query: %s", query)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	log.Printf("Retrieved user IDs: %v", userIDs)
	return userIDs, nil
}

// CreateChat creates a new chat and adds participants to it.
// It begins a transaction to insert a new chat record into the 'chats' table and associates the participants
// with the chat by inserting records into the 'chatparticipants' table.
// If any error occurs during the process, the transaction is rolled back.
func (r *PostgresChatRepository) CreateChat(ctx context.Context, chat *models.Chat, participants []int) (int, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return 0, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	var chatID int
	err = tx.QueryRow(ctx, "INSERT INTO chats (name, created_at) VALUES ($1, $2) RETURNING id", chat.Name, chat.CreatedAt).Scan(&chatID)
	if err != nil {
		log.Printf("Error inserting into chats table: %v", err)
		return 0, err
	}
	log.Printf("Chat ID: %d", chatID)

	for _, userID := range participants {
		log.Printf("Adding participant userID: %d", userID)
		_, err := tx.Exec(ctx, "INSERT INTO chatparticipants (chat_id, user_id, joined_at) VALUES ($1, $2, $3)", chatID, userID, time.Now())
		if err != nil {
			log.Printf("Error inserting into chatparticipants: %v", err)
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return 0, err
	}
	log.Println("Transaction committed successfully")

	// Возвращаем ID созданного чата
	return chatID, nil
}

// GetChatsByUserID retrieves all chats for a specific user based on their user ID.
// It queries the 'chats' table by joining it with the 'chatparticipants' table to return all chats the user is a part of.
func (r *PostgresChatRepository) GetChatsByUserID(ctx context.Context, userID int) ([]models.Chat, error) {
	query := `
        SELECT c.id, c.name, c.created_at
        FROM chats c
        JOIN chatparticipants cp ON c.id = cp.chat_id
        WHERE cp.user_id = $1`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// DeleteChat deletes a chat and all associated data, including messages and participants.
// It starts a transaction to delete all records related to the chat from the 'messages' and 'chatparticipants' tables
// and then deletes the chat itself from the 'chats' table.
// If any error occurs, the transaction is rolled back to maintain consistency.
func (r *PostgresChatRepository) DeleteChat(ctx context.Context, chatID int) error {
	// Начинаем транзакцию
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	// Обеспечиваем откат транзакции в случае ошибки
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Удаляем все сообщения чата
	_, err = tx.Exec(ctx, "DELETE FROM messages WHERE chat_id = $1", chatID)
	if err != nil {
		log.Printf("Error deleting messages for chat %d: %v", chatID, err)
		return err
	}

	// Удаляем участников чата
	_, err = tx.Exec(ctx, "DELETE FROM chatparticipants WHERE chat_id = $1", chatID)
	if err != nil {
		log.Printf("Error deleting participants for chat %d: %v", chatID, err)
		return err
	}

	// Удаляем сам чат
	_, err = tx.Exec(ctx, "DELETE FROM chats WHERE id = $1", chatID)
	if err != nil {
		log.Printf("Error deleting chat %d: %v", chatID, err)
		return err
	}

	// Коммитим транзакцию, если все операции прошли успешно
	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("Chat %d and associated data deleted successfully", chatID)
	return nil
}
