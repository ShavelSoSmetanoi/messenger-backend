package services

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/auth"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/chat"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/user"
)

type Services struct {
	User user.UserServiceInterface
	Auth auth.AuthHandlerInterface
	Chat chat.ChatServiceInterface
}

// InitServices инициализирует все сервисы и возвращает контейнер зависимостей
func InitServices() *Services {
	jw := jwtDB.NewUserTokenRepository(postgres.Db)
	ml := auth.NewAuthHandler(jw)

	// Создание репозитория пользователей
	rp := userDB.NewPostgresUserRepository(postgres.Db)
	us := user.NewUserService(rp)

	// Инициализация service chat слоя
	ch := chatDB.NewPostgresChatRepository(postgres.Db)
	chs := chat.NewChatService(ch)

	return &Services{
		User: us,
		Auth: ml,
		Chat: chs,
	}
}
