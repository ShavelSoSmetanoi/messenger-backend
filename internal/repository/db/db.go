package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

//var DB *sql.DB
//
//var ConnStr string = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
//	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
//
////func InitDB(connectionString string) (*user.PostgresUserRepository, error) {
////	//var err error
////	//DB, err = sql.Open("postgres", connectionString)
////	//if err != nil {
////	//	log.Fatalf("Failed to open database: %v", err)
////	//}
////	//
////	//// Ensure the database connection is closed when the function returns
////	//defer func() {
////	//	if err := DB.Close(); err != nil {
////	//		log.Fatalf("Failed to close database: %v", err)
////	//	}
////	//}()
////	//
////	//log.Println("Database connected successfully")
////
////	db, err := sql.Open("postgres", connectionString)
////	if err != nil {
////		log.Fatalf("Failed to open database: %v", err)
////		return nil, err
////	}
////
////	// Настройка параметров пула подключений
////	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
////	db.SetMaxIdleConns(25)                 // Максимальное количество "спящих" соединений
////	db.SetConnMaxLifetime(time.Minute * 5) // Максимальное время жизни соединения
////
////	log.Println("Database connected successfully")
////
////	return &user.PostgresUserRepository{DB: db}, nil
////}
//
//func InitDB(connectionString string) (*sql.DB, error) {
//	var err error
//	DB, err = sql.Open("postgres", connectionString)
//	if err != nil {
//		panic("pizda")
//	}
//
//	log.Println("Database connected successfully")
//
//	return DB, nil // Возвращаем соединение и nil для ошибки
//}

var db *pgxpool.Pool

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

	// Настройка параметров пула подключений
	config.MaxConns = 25                 // Максимальное количество открытых соединений
	config.MaxConnIdleTime = time.Minute // Максимальное время "ожидания" соединения

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	// Тестовое подключение
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")

	db = pool
	return db, nil
}
