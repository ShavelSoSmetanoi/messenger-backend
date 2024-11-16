package user

import (
	"context"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// ServiceInterface defines the methods that should be implemented
// by a user service responsible for handling user data and operations.
type ServiceInterface interface {
	// GetUserProfile retrieves the profile information of a user.
	// It takes the Gin context to process the request and response.
	GetUserProfile(c *gin.Context)

	// UpdateUserProfile updates the profile information of a user.
	// It takes the Gin context to handle the request and response.
	UpdateUserProfile(c *gin.Context)

	// CheckUserByUsername checks if a user exists by their username.
	// It returns a User object if found, or an error if not.
	CheckUserByUsername(username string) (*models.User, error)

	// GetUserByID retrieves a user by their unique user ID.
	// It returns the User object or an error if the user is not found.
	GetUserByID(userID string) (*models.User, error)

	// RegisterUser registers a new user into the system.
	// It processes the user registration request from the Gin context.
	RegisterUser(c *gin.Context)

	// GetSettingsByUserID retrieves the settings of a user by their user ID.
	// It returns the UserSettings object or an error if the settings are not found.
	GetSettingsByUserID(ctx context.Context, userID int) (*models.UserSettings, error)

	// UpdateSettings updates the settings of a user, such as theme and message color.
	// It takes the context and the new settings (theme and message color) for the user.
	UpdateSettings(ctx context.Context, userID int, theme, messageColor string) error

	// GetAllUsers retrieves all users in the system.
	// It returns a slice of User objects or an error if fetching users fails.
	GetAllUsers(ctx context.Context) ([]models.User, error)
}
