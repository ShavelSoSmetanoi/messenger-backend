package rest

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/user"
	middleware "github.com/ShavelSoSmetanoi/messenger-backend/internal/services/middelfare"
	user2 "github.com/ShavelSoSmetanoi/messenger-backend/internal/services/user"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRepository, err := postgres.InitDB()
	if err != nil {
		panic("Pizda")
	}

	// Создание репозитория пользователей
	rp := user.NewPostgresUserRepository(userRepository)

	// Создание сервиса пользователей
	us := user2.NewUserService(rp)

	// Маршруты для аутентификации и авторизации
	//r.POST("/login", us.)
	r.POST("/verify-email", middleware.EmailValidator())
	r.POST("/register", middleware.VerifyCode(), us.RegisterUser)

	// Проверка доступности сервиса
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	//r.GET("/profile", us.GetUserProfile)
	//authUsers := r.Group("/")
	//authUsers.Use(auth.AuthMiddleware(os.Getenv("JWT_SECRET")))
	//
	//authUsers.POST("/chats", controllers.CreateChatHandler)
	//authUsers.GET("/chats", controllers.GetChatsHandler)
	//authUsers.DELETE("/chats/:chat_id", controllers.DeleteChatHandler)
	//
	//authUsers.POST("/chats/:chat_id/messages", controllers.SendMessageHandler)
	//authUsers.GET("/chats/:chat_id/messages", controllers.GetMessagesHandler)
	//
	//authUsers.GET("/profile", controllers.GetUserProfile)
	//authUsers.PUT("/:user_id", controllers.UpdateUserProfile)
	//authUsers.GET("/check/:username", controllers.CheckUserByUsername)
	//
	//authUsers.GET("/ws", WebSoket.WebSocketHandler)
	//
	//// TODO
	//// Маршруты для работы с файлами
	////authUsers.POST("/upload", controllers.UploadFileHandler)
	////authUsers.GET("/download/:filename", controllers.DownloadFileHandler)

	return r
}
