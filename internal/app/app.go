package app

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/config"
	services2 "github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	"log"
)

func Run() {
	config.SetupConfig()

	// Инициализация сервисов
	services := services2.InitServices()

	// Создание роутера и инициализация обработчиков
	handler := rest.NewHandler(services)
	router := handler.Init()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
