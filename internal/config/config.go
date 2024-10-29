package config

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	rd "github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/redis"
	"github.com/joho/godotenv"
	"log"
)

func SetapConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rd.InitRedis()

	_, err1 := postgres.InitDB()
	if err1 != nil {
		panic("Pizda")
	}
}
