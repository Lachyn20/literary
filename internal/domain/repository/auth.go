package repository

import (
	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenGenerator interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID uuid.UUID
	Role   string
}

func NewTokenClaims(userID uuid.UUID, role string) *TokenClaims {
	return &TokenClaims{UserID: userID, Role: role}
}
