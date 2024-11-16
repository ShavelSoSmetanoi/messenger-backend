package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// Service struct contains the user repository, which is used for user data access.
type Service struct {
	userRepo userDB.UserRepository
}

// NewUserService creates and returns a new Service instance with the provided user repository.
func NewUserService(repo userDB.UserRepository) *Service {
	return &Service{
		userRepo: repo,
	}
}

// GetSettingsByUserID retrieves the settings for a user based on their user ID.
func (h *Service) GetSettingsByUserID(ctx context.Context, userID int) (*models.UserSettings, error) {
	return h.userRepo.GetSettingsByUserID(ctx, userID)
}

// UpdateSettings updates the user's theme and message color preferences.
func (h *Service) UpdateSettings(ctx context.Context, userID int, theme, messageColor string) error {
	return h.userRepo.UpdateSettings(ctx, userID, theme, messageColor)
}

// GetUserByID retrieves a user by their user ID.
func (h *Service) GetUserByID(userID string) (*models.User, error) {
	user, err := h.userRepo.GetUserByID(context.Background(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error fetching user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

// RegisterUser handles the user registration process.
func (h *Service) RegisterUser(c *gin.Context) {
	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User data not found.",
		})
		return
	}

	data := userData.(map[string]string)

	username := data["username"]
	email := data["email"]
	password := data["password"]

	err := h.userRepo.CreateUser(username, email, password, "", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// GetUserProfile retrieves the profile of a user based on their user ID.
func (h *Service) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, err := h.userRepo.GetUserByID(context.Background(), userID.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"uuid":     user.UniqueId,
		"username": user.Username,
		"email":    user.Email,
		"photo":    user.Photo,
		"about":    user.About,
	})
}

// UpdateUserProfile updates the profile information of a user based on input from the request body.
func (h *Service) UpdateUserProfile(c *gin.Context) {
	userID := c.Param("user_id")
	var userUpdate models.UserUpdate

	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userRepo.UpdateUser(context.Background(), userID, userUpdate)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// CheckUserByUsername checks if a user exists with the given username.
func (h *Service) CheckUserByUsername(username string) (*models.User, error) {
	user, err := h.userRepo.GetUserByUsername(context.Background(), username)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves a list of all users from the database.
func (h *Service) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := h.userRepo.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		return nil, err
	}
	return users, nil
}
