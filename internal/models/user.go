package models

// User represents a user in the system.
type User struct {
	ID       string `json:"id"`        // Unique identifier for the user
	Username string `json:"username"`  // Username chosen by the user
	Email    string `json:"email"`     // Email address of the user
	Password string `json:"password"`  // Encrypted password for the user
	Photo    []byte `json:"photo"`     // Profile photo of the user (as a byte array)
	UniqueId string `json:"unique_id"` // Unique ID for the user (e.g., UUID)
	About    string `json:"about"`     // Short description or biography of the user
}

// UserUpdate represents the fields that can be updated in a user's profile.
type UserUpdate struct {
	About string `json:"about"` // Updated description or biography of the user
	Photo []byte `json:"photo"` // Updated profile photo of the user (as a byte array)
}

// UserResponse is the format used for sending user data in API responses.
type UserResponse struct {
	ID       string `json:"id"`       // Unique identifier for the user
	Username string `json:"username"` // Username of the user
	Email    string `json:"email"`    // Email address of the user
	Photo    []byte `json:"photo"`    // Profile photo of the user (as a byte array)
	About    string `json:"about"`    // Description or biography of the user
}
