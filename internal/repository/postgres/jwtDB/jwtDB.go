package jwtDB

import (
	"context"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	//"runtime/debug"
	"time"
)

// UserTokenRepositoryInterface defines the methods for working with user tokens
type UserTokenRepositoryInterface interface {
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	SaveToken(ctx context.Context, userID, token string) error
	DeleteToken(ctx context.Context, token string) error
	IsTokenValid(ctx context.Context, token string) (bool, error)
	GetTokensByUserID(ctx context.Context, userID string) ([]models.UserToken, error)
}

// UserTokenRepository struct implements UserTokenRepositoryInterface
type UserTokenRepository struct {
	DB *pgxpool.Pool
}

// NewUserTokenRepository creates and returns a new UserTokenRepository instance
func NewUserTokenRepository(db *pgxpool.Pool) *UserTokenRepository {
	return &UserTokenRepository{DB: db}
}

// AuthenticateUser checks if the provided username and password are valid
// If valid, it returns the user's information; otherwise, an error is returned
func (r *UserTokenRepository) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User

	err := r.DB.QueryRow(ctx, "SELECT id, username, email, password, photo, unique_id, about FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// SaveToken saves a user token in the database
func (r *UserTokenRepository) SaveToken(ctx context.Context, userID, token string) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO user_tokens (user_id, token, created_at) VALUES ($1, $2, $3)",
		userID, token, time.Now())
	return err
}

// DeleteToken deletes a user token from the database based on the token value
func (r *UserTokenRepository) DeleteToken(ctx context.Context, token string) error {
	_, err := r.DB.Exec(ctx, "DELETE FROM user_tokens WHERE token = $1", token)
	return err
}

// IsTokenValid checks if the provided token is valid and if it has expired
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

// GetTokensByUserID retrieves all tokens associated with a specific user
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
