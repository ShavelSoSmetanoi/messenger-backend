package userDB

import (
	"context"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// UserRepository defines the interface for interacting with user data.
type UserRepository interface {
	CreateUser(username string, email string, password string, about string, photo []byte) error
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetSettingsByUserID(ctx context.Context, userID int) (*models.UserSettings, error)
	UpdateSettings(ctx context.Context, userID int, theme, messageColor string) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

// PostgresUserRepository is a concrete implementation of UserRepository for PostgresSQL.
type PostgresUserRepository struct {
	DB *pgxpool.Pool
}

// NewPostgresUserRepository creates a new instance of PostgresUserRepository.
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

// GetSettingsByUserID retrieves user settings by ID from the user_settings table. If no settings are found, it returns an error.
func (r *PostgresUserRepository) GetSettingsByUserID(ctx context.Context, userID int) (*models.UserSettings, error) {
	var settings models.UserSettings

	query := `SELECT id, user_id, theme, message_color, created_at FROM user_settings WHERE user_id = $1`
	err := r.DB.QueryRow(ctx, query, userID).
		Scan(&settings.ID, &settings.UserID, &settings.Theme, &settings.MessageColor, &settings.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No settings found for user ID: %d", userID)
			return nil, errors.New("settings not found")
		}
		log.Printf("Error querying user settings: %v", err)
		return nil, err
	}

	return &settings, nil
}

// UpdateSettings обновляет тему и цвет сообщений пользователя.
func (r *PostgresUserRepository) UpdateSettings(ctx context.Context, userID int, theme, messageColor string) error {
	query := `UPDATE user_settings SET theme = COALESCE(NULLIF($1, ''), theme), message_color = COALESCE(NULLIF($2, ''), message_color) WHERE user_id = $3`
	_, err := r.DB.Exec(ctx, query, theme, messageColor, userID)
	if err != nil {
		log.Printf("Error updating settings for user ID %d: %v", userID, err)
		return err
	}

	log.Printf("Settings for user ID %d updated successfully", userID)
	return nil
}

// CreateUser creates a new user with the provided details, hashes the password,
// generates a unique ID, and inserts the user into the database.
func (r *PostgresUserRepository) CreateUser(username, email, password, about string, photo []byte) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	uniqueID := pkg.GenerateUniqueID()

	query := `INSERT INTO users (username, email, password, photo, unique_id, about, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err = r.DB.Exec(context.Background(), query, username, email, hashedPassword, photo, uniqueID, about)
	if err != nil {
		log.Printf("Error inserting user into database: %v, Username: %s, Email: %s, About: %s", err, username, email, about)
		return err
	}

	return nil
}

// AuthenticateUser verifies the provided credentials for a user with the given username.
func (r *PostgresUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, email, password, photo, unique_id, about 
	          FROM users WHERE username = $1`
	err := r.DB.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No user found with username: %s", username)
			return nil, errors.New("invalid credentials")
		}
		log.Printf("Error querying user: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", username)
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// GetUserByID retrieves a user's details by their ID. Returns an error if the user is not found.
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, email, photo, unique_id, about 
	          FROM users WHERE id = $1`
	err := r.DB.QueryRow(ctx, query, userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.Photo, &user.UniqueId, &user.About)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No user found with ID: %s", userID)
			return nil, errors.New("user not found")
		}
		// Логирование ошибки запроса
		log.Printf("Error querying user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates the specified fields of a user in the database using the provided userID and update data.
// Returns an error if the update fails.
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error {
	var photo interface{}
	if len(userUpdate.Photo) == 0 {
		photo = nil
	} else {
		photo = userUpdate.Photo
	}

	query := `
		UPDATE users 
		SET about = $1, photo = COALESCE($2::bytea, photo)
		WHERE id = $3`

	_, err := r.DB.Exec(ctx, query, userUpdate.About, photo, userID)
	if err != nil {
		log.Printf("Error updating user with ID %s: %v", userID, err)
		return err
	}

	return nil
}

// GetUserByUsername retrieves a user by their username. Returns nil if the user does not exist.
func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, email, password, photo, unique_id, about 
	          FROM users WHERE username = $1`
	err := r.DB.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No user found with username: %s", username)
			return nil, nil
		}
		log.Printf("Error querying user by username: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *PostgresUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, photo, about FROM users`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Photo, &user.About); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
