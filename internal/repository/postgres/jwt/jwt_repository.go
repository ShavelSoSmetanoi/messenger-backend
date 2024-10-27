package jwt

import (
	"context"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserTokenRepository struct {
	DB *pgxpool.Pool
}

func NewUserTokenRepository(db *pgxpool.Pool) *UserTokenRepository {
	return &UserTokenRepository{DB: db}
}

// SaveToken сохраняет токен пользователя в базе данных
func (r *UserTokenRepository) SaveToken(ctx context.Context, userID, token string) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO user_tokens (user_id, token, created_at) VALUES ($1, $2, $3)",
		userID, token, time.Now())
	return err
}

// DeleteToken удаляет токен из базы данных
func (r *UserTokenRepository) DeleteToken(ctx context.Context, token string) error {
	_, err := r.DB.Exec(ctx, "DELETE FROM user_tokens WHERE token = $1", token)
	return err
}

// IsTokenValid проверяет, действителен ли токен, и истек ли срок его действия
func (r *UserTokenRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	var createdAt time.Time
	err := r.DB.QueryRow(ctx, "SELECT created_at FROM user_tokens WHERE token = $1", token).Scan(&createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, errors.New("token not found")
		}
		return false, err
	}

	if time.Since(createdAt) > 24*time.Hour {
		return false, errors.New("token expired")
	}

	return true, nil
}

// GetTokensByUserID возвращает все токены пользователя
func (r *UserTokenRepository) GetTokensByUserID(ctx context.Context, userID string) ([]models.UserToken, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, user_id, token, created_at FROM user_tokens WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []models.UserToken
	for rows.Next() {
		var token models.UserToken
		if err := rows.Scan(&token.ID, &token.UserID, &token.Token, &token.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}
