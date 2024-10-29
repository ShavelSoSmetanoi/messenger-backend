package main

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/app"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/rest"
	_ "github.com/lib/pq"
)

func main() {
	// Запускаем сервис
	app.Run()
}
