package users

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSession(userID string) *Session {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	sessionToken := hex.EncodeToString(b)

	return &Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     sessionToken,
		CreatedAt: time.Now(),
	}
}
