package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

var ConnStr string = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

func InitDB(connectionString string) {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Ensure the database connection is closed when the function returns
	defer func() {
		if err := DB.Close(); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()

	log.Println("Database connected successfully")
}
