package main

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	//defer repository.Db.Close()
	//transport.InitRedis()

	//db.InitDB(db.ConnStr)
	//Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Настройка маршрутизатора
	router := rest.SetupRouter()

	// Запуск HTTP-сервера
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	//// Создание нового приложения
	//myApp := app.NewApp()
	//
	//// Запуск приложения на порту 8080 (или другом, который вы хотите использовать)
	//myApp.Run("8080")
}
