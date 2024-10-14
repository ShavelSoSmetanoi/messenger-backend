package pkg

import (
	"crypto/rand"
	"fmt"
)

// Функция для генерации случайного 6-значного кода
func GenerateCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%06d", b)
}
