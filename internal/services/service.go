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
	Auth    auth.Interface
	Chat    chat.ServiceInterface
	Message message.ServiceInterface
	File    file.Service
}

// InitServices initializes all services and returns a dependency container
func InitServices() *Services {
	// Initialize the JWT token repository and authentication handler
	jw := jwtDB.NewUserTokenRepository(postgres.Db)
	ml := auth.NewAuthHandler(jw)

	// Create the user repository and user service
	rp := userDB.NewPostgresUserRepository(postgres.Db)
	us := user.NewUserService(rp)

	// Initialize message repository and message service
	ms := messageDB.NewPostgresMessageRepository(postgres.Db)
	ckl := chatparticipantDB.NewPostgresChatParticipantRepository(postgres.Db)
	mss := message.NewMessageService(ms, ckl)

	// Initialize the chat repository and chat service
	ch := chatDB.NewPostgresChatRepository(postgres.Db)
	chs := chat.NewChatService(ch, ckl)

	// Initialize the S3 client for file storage
	s3Client, err := s3.NewS3Client(os.Getenv("BUCKET_NAME"))
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	fileService := file.NewS3FileService(s3Client.Client, s3Client.Bucket, ms)

	// Return the Services struct containing all initialized services
	return &Services{
		User:    us,
		Auth:    ml,
		Chat:    chs,
		Message: mss,
		File:    fileService,
	}
}
