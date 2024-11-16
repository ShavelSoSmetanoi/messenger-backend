package config

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	rd "github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/redis"
	"github.com/joho/godotenv"
	"log"
	"time"
)

// PostgresConfig holds the configuration details for connecting to a PostgresSQL database.
type PostgresConfig struct {
	User           string
	Password       string
	Host           string
	DBName         string
	MaxConnections int32
	IdleTimeout    time.Duration
}

// RedisConfig holds the configuration details for connecting to a Redis server.
type RedisConfig struct {
	Address string
	DB      int
	Timeout time.Duration
}

// SetupConfig initializes the configuration by loading environment variables and initializing the services.
func SetupConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Redis connection
	rd.InitRedis()

	// Initialize PostgresSQL connection and check for errors
	_, err1 := postgres.InitDB()
	if err1 != nil {
		// If the database connection fails, panic with a message
		panic("db panic")
	}
}
