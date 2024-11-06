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

	r.GET("/user/:user_id", h.GetUserByIDHandler)
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

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Photo    []byte `json:"photo"`
	About    string `json:"about"`
}

func (h *Handler) GetUserByIDHandler(c *gin.Context) {
	// Получаем user_id из параметров маршрута
	userID := c.Param("user_id")

	// Вызов метода бизнес-логики для получения пользователя по ID
	user, err := h.services.User.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Photo:    user.Photo,
		About:    user.About,
	}

	c.JSON(http.StatusOK, userResponse)
}
