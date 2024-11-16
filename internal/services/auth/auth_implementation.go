package auth

import (
	"context"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg/JWT"
)

// Handler AuthHandler handles authentication-related logic.
type Handler struct {
	authService jwtDB.UserTokenRepositoryInterface
}

// NewAuthHandler creates a new AuthHandler with the provided authentication service.
func NewAuthHandler(repo jwtDB.UserTokenRepositoryInterface) *Handler {
	return &Handler{
		authService: repo,
	}
}

// Login authenticates the user, validates existing tokens, and generates a new token if necessary.
func (h *Handler) Login(username, password string) (string, error) {
	// Authenticate the user using the provided username and password
	user, err := h.authService.AuthenticateUser(context.Background(), username, password)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Retrieve existing tokens associated with the user
	tokens, err := h.authService.GetTokensByUserID(context.Background(), user.ID)
	if err != nil {
		return "", errors.New("failed to retrieve tokens")
	}

	// Check if any of the tokens are valid
	var validToken string
	for _, t := range tokens {
		valid, err := h.authService.IsTokenValid(context.Background(), t.Token)
		if err == nil && valid {
			validToken = t.Token
			break
		}
	}

	// Generate a new token if no valid tokens are found
	if validToken == "" {
		token, err := JWT.CreateJWT(user.ID)
		if err != nil {
			return "", errors.New("failed to generate token")
		}
		// Save the newly generated token in the repository
		if err := h.authService.SaveToken(context.Background(), user.ID, token); err != nil {
			return "", errors.New("failed to save token")
		}
		validToken = token
	}

	// Return the valid token (either existing or newly created)
	return validToken, nil
}
