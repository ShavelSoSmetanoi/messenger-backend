package postgres

import "github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/postgres/userDB"

type Deps struct {
	User userDB.UserRepository
}
