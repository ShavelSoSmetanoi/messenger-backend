package rest

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	middleware "github.com/ShavelSoSmetanoi/messenger-backend/internal/services/middelfare"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/Websocket"
	"github.com/gin-gonic/gin"
	"os"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	r := gin.Default()

	r.Group("/")
	{
		h.InitAuthRouter(r)
	}

	r.GET("/ws", Websocket.WebSocketHandler)

	authUsers := r.Group("/")
	authUsers.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))

	h.InitUserRouter(authUsers)
	h.InitChatRouter(authUsers)
	h.InitMessageRouter(authUsers)

	return r
}
