package userDB

import (
	"context"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository interface {
	CreateUser(username string, email string, password string, about string, photo []byte) error
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	//GetUserByID(ctx context.Context, userID string) (*models.User, error)
	//UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error
	//GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type PostgresUserRepository struct {
	DB *pgxpool.Pool
}

func (r *PostgresUserRepository) GetUserProfile(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresUserRepository) UpdateUserProfile(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresUserRepository) CheckUserByUsername(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresUserRepository) RegisterUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

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

// autch
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

//func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
//	var user *models.User
//
//	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, photo, unique_id, about FROM users WHERE id = $1", userID).
//		Scan(&user.ID, &user.Username, &user.Email, &user.Photo, &user.UniqueId, &user.About)
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}
//
//func (r *PostgresUserRepository) UpdateUser(ctx context.Context, userID string, userUpdate models.UserUpdate) error {
//	_, err := r.db.ExecContext(ctx, "UPDATE users SET email = $1, about = $2, photo = $3 WHERE id = $4",
//		userUpdate.Email, userUpdate.About, userUpdate.Photo, userID)
//	return err
//}
//
//func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
//	var user *models.User
//	query := "SELECT id, username, email, password, photo, unique_id, about FROM users WHERE username = $1"
//
//	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Photo, &user.UniqueId, &user.About)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil // Пользователь не найден
//		}
//		return nil, fmt.Errorf("error getting user by username: %w", err)
//	}
//
//	return user, nil
//}
