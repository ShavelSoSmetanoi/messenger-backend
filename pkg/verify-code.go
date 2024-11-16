package pkg

import (
	"crypto/rand"
	"fmt"
	"log"
)

// GenerateCode generates a random 6-digit code as a string.
// It uses a cryptographic random number generator to ensure secure random values.
func GenerateCode() string {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating random code: %v", err)
		return ""
	}
	return fmt.Sprintf("%06d", b)
}
