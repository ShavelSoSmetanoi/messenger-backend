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

type UserRepository interface {
	CreateUser(username string, email string, password string, about string, photo []byte) error
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type PostgresUserRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

// Создание пользователя
func (r *PostgresUserRepository) CreateUser(username, email, password, about string, photo []byte) error {
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	// Генерация уникального ID
	uniqueID := pkg.GenerateUniqueID()

	// Выполнение запроса на вставку данных пользователя
	query := `INSERT INTO users (username, email, password, photo, unique_id, about, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err = r.DB.Exec(context.Background(), query, username, email, hashedPassword, photo, uniqueID, about)
	if err != nil {
		log.Printf("Error inserting user into database: %v, Username: %s, Email: %s, About: %s", err, username, email, about)
		return err
	}

	// Логирование успешного выполнения
	log.Printf("User %s created successfully", username)
	return nil
}

// Проверка аунтификакации пользователя
func (r *PostgresUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User

	// Выполнение запроса на получение данных пользователя
	query := `SELECT id, username, email, password, photo, unique_id, about 
	          FROM users WHERE username = $1`
	err := r.DB.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)

	if err != nil {
		// Если пользователь не найден
		if err == pgx.ErrNoRows {
			log.Printf("No user found with username: %s", username)
			return nil, errors.New("invalid credentials")
		}
		// Логирование ошибки запроса
		log.Printf("Error querying user: %v", err)
		return nil, err
	}

	// Сравнение хешированного пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", username)
		return nil, errors.New("invalid credentials")
	}

	// Логирование успешной аутентификации
	log.Printf("User %s authenticated successfully", username)
	return &user, nil
}

// Получение пользователя по ID
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User

	// Выполнение запроса на получение данных пользователя
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

// Обновление пользователя
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error {
	// Выполнение запроса на обновление данных пользователя
	query := `UPDATE users SET email = $1, about = $2, photo = $3 WHERE id = $4`
	_, err := r.DB.Exec(ctx, query, userUpdate.About, userUpdate.Photo, userID)

	if err != nil {
		log.Printf("Error updating user with ID %s: %v", userID, err)
		return err
	}

	// Логирование успешного обновления
	log.Printf("User with ID %s updated successfully", userID)
	return nil
}

// Получение пользователя по имени пользователя
func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	// Выполнение запроса на получение данных пользователя
	query := `SELECT id, username, email, password, photo, unique_id, about 
	          FROM users WHERE username = $1`
	err := r.DB.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No user found with username: %s", username)
			return nil, nil // Пользователь не найден
		}
		// Логирование ошибки запроса
		log.Printf("Error querying user by username: %v", err)
		return nil, err
	}

	return &user, nil
}
