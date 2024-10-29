package rest

import (
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/services"
	"github.com/gin-gonic/gin"
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
		h.InitSetupRouter(r)
	}

	return r
}
