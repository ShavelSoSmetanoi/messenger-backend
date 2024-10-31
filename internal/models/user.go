package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    []byte `json:"photo"`
	UniqueId string `json:"unique_id"` // @гусь000000000000
	About    string `json:"about"`
}

type UserUpdate struct {
	Email string `json:"email"` // - delete
	About string `json:"about"`
	Photo []byte `json:"photo"`
}
