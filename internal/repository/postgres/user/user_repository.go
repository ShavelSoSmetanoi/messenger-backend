package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type UserRepository interface {
	//CreateUser(username string, email string, password string, about string, photo []byte) error
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Close закрывает соединение с базой данных
func (r *PostgresUserRepository) Close() error {
	return r.db.Close()
}

func CreateUser(username, email, password, about string, photo []byte) error {
	//Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	// Строка подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&client_encoding=UTF8",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}
	defer db.Close()

	fmt.Println("Successfully connected to PostgreSQL!")

	//Генерация уникального ID
	uniqueID := pkg.GenerateUniqueID()

	// Выполнение запроса на вставку данных пользователя
	_, err = db.Exec("INSERT INTO users (username, email, password, photo, unique_id, about, registration_date) VALUES ($1, $2, $3, $4, $5, $6, NOW())",
		username, email, hashedPassword, photo, uniqueID, about)
	if err != nil {
		log.Printf("Error inserting user into database: %v, Username: %s, Email: %s, About: %s", err, username, email, about)
		return err
	}

	// Логирование успешного выполнения
	log.Printf("User %s created successfully", username)
	return nil
}

func (r *PostgresUserRepository) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	var user *models.User

	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, password, photo, unique_id, about FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)
	if err != nil {
		if err == sql.ErrNoRows {
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

	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user *models.User

	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, photo, unique_id, about FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.Photo, &user.UniqueId, &user.About)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET email = $1, about = $2, photo = $3 WHERE id = $4",
		userUpdate.Email, userUpdate.About, userUpdate.Photo, userID)
	return err
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user *models.User
	query := "SELECT id, username, email, password, photo, unique_id, about FROM users WHERE username = $1"

	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}

	return user, nil
}
