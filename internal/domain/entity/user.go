package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
