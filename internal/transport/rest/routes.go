package rest

import (
	"database/sql"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/user"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	var connStr string = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
	}
	defer db.Close()

	rp := user.NewPostgresUserRepository(db)
	us := services.NewUserService(rp)

	// Маршруты для аутентификации и авторизации
	//r.POST("/login", us.)
	r.POST("/register", us.RegisterUser)

	// Проверка доступности сервиса
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/profile", us.GetUserProfile)
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
