package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    []byte `json:"photo"`
	UniqueId string `json:"unique_id"`
	About    string `json:"about"`
}

type UserUpdate struct {
	About string `json:"about"`
	Photo []byte `json:"photo"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Photo    []byte `json:"photo"`
	About    string `json:"about"`
}
