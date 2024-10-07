package main

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/app"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
)

func main() {
	//r := rest.SetupRouter()
	//
	//// Запускаем сервер на порту 8080
	//r.Run(":8080")

	// Создание нового приложения
	myApp := app.NewApp()

	// Запуск приложения на порту 8080 (или другом, который вы хотите использовать)
	myApp.Run("8080")
}
