package chat

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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

func (r *PostgresChatRepository) GetUserIDsByNicknames(ctx context.Context, nicknames []string) ([]int, error) {
	// TODO: Implement logic to retrieve user IDs based on nicknames
	log.Println("GetUserIDsByNicknames not implemented yet")
	return nil, nil
}

func (r *PostgresChatRepository) CreateChat(ctx context.Context, chat *models.Chat, participants []int) error {
	// TODO: Implement logic to create a chat and add participants
	log.Println("CreateChat not implemented yet")
	return nil
}

func (r *PostgresChatRepository) GetChatsByUserID(ctx context.Context, userID int) ([]models.Chat, error) {
	// TODO: Implement logic to retrieve all chats for a specific user
	log.Println("GetChatsByUserID not implemented yet")
	return nil, nil
}
