package auth

import (
	"context"
	"errors"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg/JWT"
)

type AuthHandler struct {
	authService jwtDB.UserTokenRepositoryInterface
}

func NewAuthHandler(repo jwtDB.UserTokenRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		authService: repo,
	}
}

// Login authenticates the user, manages tokens, and returns a valid token or an error.
func (h *AuthHandler) Login(username, password string) (string, error) {
	// Аутентификация пользователя (предполагается, что models.AuthenticateUser возвращает user.ID при успешной проверке)
	user, err := h.authService.AuthenticateUser(context.Background(), username, password)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Получение токенов, связанных с пользователем
	tokens, err := h.authService.GetTokensByUserID(context.Background(), user.ID)
	if err != nil {
		return "", errors.New("failed to retrieve tokens")
	}

	// Проверка валидности токенов
	var validToken string
	for _, t := range tokens {
		valid, err := h.authService.IsTokenValid(context.Background(), t.Token)
		if err == nil && valid {
			validToken = t.Token
			break
		}
	}

	// Создание нового токена, если нет действующего
	if validToken == "" {
		token, err := JWT.CreateJWT(user.ID)
		if err != nil {
			return "", errors.New("failed to generate token")
		}
		// Сохранение нового токена
		if err := h.authService.SaveToken(context.Background(), user.ID, token); err != nil {
			return "", errors.New("failed to save token")
		}
		validToken = token
	}

	// Возврат действующего или вновь созданного токена
	return validToken, nil
}
