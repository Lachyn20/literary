package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{Cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	return string(bytes), err
}

func (h *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
