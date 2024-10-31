package auth

// AuthHandlerInterface defines the methods for handling authentication-related requests
type AuthHandlerInterface interface {
	Login(username, password string) (string, error)
}
