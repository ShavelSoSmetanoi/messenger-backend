package JWT

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

// CreateJWT generates a JWT for the given user ID and returns it as a string.
// The token is signed using the HS256 algorithm.
func CreateJWT(userID string) (string, error) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	token := jwt.New(jwt.SigningMethodHS256)

	// Установка клеймов (payload) токена
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
