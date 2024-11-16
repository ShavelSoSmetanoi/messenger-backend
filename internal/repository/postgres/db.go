package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var Db *pgxpool.Pool

// InitDB initializes a connection to the PostgresSQL database and returns the connection pool.
func InitDB() (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
		return nil, err
	}

	// Set connection pool parameters
	config.MaxConns = 25                     // Maximum number of open connections
	config.MaxConnIdleTime = time.Minute * 1 // Maximum idle time for a connection

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")

	Db = pool
	return Db, nil
}
