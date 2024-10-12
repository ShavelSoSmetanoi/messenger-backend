package app

import (
	"database/sql"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/user"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

// App - структура для основного приложения
type App struct {
	router        *gin.Engine
	userService   services.UserServiceInterface
	userTransport *rest.UserTransport
}

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

// NewApp - создает новое приложение
func NewApp() *App {
	// Инициализация роутера
	router := gin.Default()

	InitDB(ConnStr)
	// Инициализация репозиториев
	userRepo := user.NewPostgresUserRepository(DB)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo)

	// Инициализация транспортных слоев
	userTransport := rest.NewUserTransport(userService)

	// Создание нового приложения
	app := &App{
		router:        router,
		userService:   userService,
		userTransport: userTransport,
	}

	// Регистрируем маршруты
	app.registerRoutes()

	return app
}

// registerRoutes - регистрирует маршруты приложения
func (a *App) registerRoutes() {
	a.userTransport.RegisterRoutes(a.router)
}

// Run - запускает приложение
func (a *App) Run(port string) {
	if err := a.router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
