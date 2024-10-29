package config

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	rd "github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/redis"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type PostgresConfig struct {
	User           string
	Password       string
	Host           string
	DBName         string
	MaxConnections int32
	IdleTimeout    time.Duration
}

type RedisConfig struct {
	Address string
	DB      int
	Timeout time.Duration
}

func SetupConfig() {
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
