package models

import (
	"time"
)

type UserStore interface {
	GetUserByMail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user RegUser) error
	GetUsers() (*[]PublicUser, error)
}

type RegUser struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type PublicUser struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}
