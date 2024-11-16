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

// NewHandler creates and returns a new Handler with the provided services.
func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services: services,
	}
}

// Init initializes the router and sets up all routes
// It returns a Gin engine with all configured routes
func (h *Handler) Init() *gin.Engine {
	r := gin.Default()

	r.Group("/")
	{
		h.InitAuthRouter(r)
	}

	r.GET("/ws", Websocket.Handler)

	authUsers := r.Group("/")
	authUsers.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))

	h.InitUserRouter(authUsers)
	h.InitChatRouter(authUsers)
	h.InitMessageRouter(authUsers)
	h.InitFileRoutes(authUsers)

	return r
}
