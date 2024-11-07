package services

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/chatparticipantDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/jwtDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/messageDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/s3"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/auth"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/chat"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/file"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/message"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services/user"
	"log"
	"os"
)

type Services struct {
	User    user.ServiceInterface
	Auth    auth.AuthHandlerInterface
	Chat    chat.ServiceInterface
	Message message.ServiceInterface
	File    file.FileService
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

	ms := messageDB.NewPostgresMessageRepository(postgres.Db)
	ckl := chatparticipantDB.NewPostgresChatParticipantRepository(postgres.Db)
	mss := message.NewMessageService(ms, ckl)

	s3Client, err := s3.NewS3Client(os.Getenv("BUCKET_NAME"))
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	fileService := file.NewS3FileService(s3Client.Client, s3Client.Bucket)

	return &Services{
		User:    us,
		Auth:    ml,
		Chat:    chs,
		Message: mss,
		File:    fileService,
	}
}
