package rest

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitUserRouter(r *gin.RouterGroup) {

	r.GET("/profile", h.services.User.GetUserProfile)

	r.PUT("/:user_id", h.services.User.UpdateUserProfile)

	r.GET("/check/:username", h.services.User.CheckUserByUsername)
}
