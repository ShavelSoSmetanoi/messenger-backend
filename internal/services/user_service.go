package services

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/models"

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
