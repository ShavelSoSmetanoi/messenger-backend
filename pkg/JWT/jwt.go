package JWT

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func CreateJWT(userID string) (string, error) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	// Создание нового токена с использованием HS256 алгоритма подписи
	token := jwt.New(jwt.SigningMethodHS256)

	// Установка клеймов (payload) токена
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Токен действителен в течение 1 часа

	// Подписание токена с использованием секретного ключа
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
