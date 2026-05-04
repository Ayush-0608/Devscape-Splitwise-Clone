package models

type User struct {
	ID       int64  `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
