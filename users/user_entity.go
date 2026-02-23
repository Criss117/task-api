package users

import (
	"time"

	"github.com/google/uuid"
)

type SignUpDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserWithSession struct {
	User
	Session *Session `json:"session"`
}

func NewUser(name string, email string, password string) *User {
	userId := uuid.NewString()

	return &User{
		ID:        userId,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
