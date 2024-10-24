package pkg

import (
	"crypto/rand"
	"fmt"
	"log"
)

// GenerateCode - Функция для генерации случайного 6-значного кода
func GenerateCode() string {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating random code: %v", err)
		return ""
	}
	return fmt.Sprintf("%06d", b)
}
