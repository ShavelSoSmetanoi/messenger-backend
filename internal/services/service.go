package services

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/user"
)

type Services struct {
	User user.UserServiceInterface
}

// InitServices инициализирует все сервисы и возвращает контейнер зависимостей
func InitServices() *Services {
	// Создание репозитория пользователей
	rp := userDB.NewPostgresUserRepository(postgres.Db)
	us := user.NewUserService(rp)

	return &Services{
		User: us,
	}
}
