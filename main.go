package main

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/app"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	_ "github.com/lib/pq"
)

func main() {
	//Hello
	//Загрузка переменных окружения из файла .env

	app.Run()

	//// Создание нового приложения
	//myApp := app.NewApp()
	//
	//// Запуск приложения на порту 8080 (или другом, который вы хотите использовать)
	//myApp.Run("8080")
}
