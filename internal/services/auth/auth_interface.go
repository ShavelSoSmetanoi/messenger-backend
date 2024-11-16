package auth

// Interface defines the methods for handling authentication-related requests
type Interface interface {
	Login(username, password string) (string, error)
}
