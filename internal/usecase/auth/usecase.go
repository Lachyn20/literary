package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type RegisterUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher repository.PasswordHasher
	tokenGen       repository.TokenGenerator
	tokenRepo      repository.RefreshTokenRepository
}

func NewRegisterUseCase(userRepo repository.UserRepository, passwordHasher repository.PasswordHasher, tokenGen repository.TokenGenerator, tokenRepo repository.RefreshTokenRepository) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo, passwordHasher: passwordHasher, tokenGen: tokenGen, tokenRepo: tokenRepo}
}

func (u *RegisterUseCase) Execute(ctx context.Context, name, email, password, role string) (*entity.User, string, string, error) {
	_, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, "", "", errors.New("email already registered")
	}

	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return nil, "", "", err
	}

	user := &entity.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, "", "", err
	}

	accessToken, err := u.tokenGen.GenerateAccessToken(user.ID, role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := u.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", err
	}

	if err := u.tokenRepo.Create(ctx, &entity.RefreshToken{ID: uuid.New(), UserID: user.ID, Token: refreshToken, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), CreatedAt: time.Now()}); err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
