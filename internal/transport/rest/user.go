package rest

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) InitUserRouter(r *gin.RouterGroup) {

	r.GET("/profile", h.services.User.GetUserProfile)

	r.PUT("/:user_id", h.services.User.UpdateUserProfile)

	r.GET("/check/:username", h.CheckUserByUsernameHandler)
}

func (h *Handler) CheckUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")

	user, err := h.services.User.CheckUserByUsername(username) // Вызов метода бизнес-логики
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
