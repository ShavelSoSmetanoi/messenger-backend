package services

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	_ "github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/auth"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/user"
)

type Services struct {
	User user.UserServiceInterface
	Auth auth.AuthHandlerInterface
}

// InitServices инициализирует все сервисы и возвращает контейнер зависимостей
func InitServices() *Services {
	jw := jwtDB.NewUserTokenRepository(postgres.Db)
	ml := auth.NewAuthHandler(*jw)

	// Создание репозитория пользователей
	rp := userDB.NewPostgresUserRepository(postgres.Db)
	us := user.NewUserService(rp)

	return &Services{
		User: us,
		Auth: ml,
	}
}
