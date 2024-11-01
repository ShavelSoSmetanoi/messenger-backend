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

type ChatRepository interface {
	GetUserIDsByNicknames(ctx context.Context, nicknames []string) ([]int, error)
	CreateChat(ctx context.Context, chat *models.Chat, participants []int) error
	GetChatsByUserID(ctx context.Context, userID int) ([]models.Chat, error)
}

type PostgresChatRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresChatRepository(db *pgxpool.Pool) *PostgresChatRepository {
	return &PostgresChatRepository{DB: db}
}

// TODO - проверить!
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

// CreateChat создает новый чат и добавляет участников
func (r *PostgresChatRepository) CreateChat(ctx context.Context, chat *models.Chat, participants []int) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
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
		return err
	}
	log.Printf("Chat ID: %d", chatID)

	for _, userID := range participants {
		log.Printf("Adding participant userID: %d", userID)
		_, err := tx.Exec(ctx, "INSERT INTO chatparticipants (chat_id, user_id, joined_at) VALUES ($1, $2, $3)", chatID, userID, time.Now())
		if err != nil {
			log.Printf("Error inserting into chatparticipants: %v", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	log.Println("Transaction committed successfully")
	return nil
}

// GetChatsByUserID получает все чаты для конкретного пользователя
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
