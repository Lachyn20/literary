package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type LoginUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher repository.PasswordHasher
	tokenGen       repository.TokenGenerator
	tokenRepo      repository.RefreshTokenRepository
}

func NewLoginUseCase(userRepo repository.UserRepository, passwordHasher repository.PasswordHasher, tokenGen repository.TokenGenerator, tokenRepo repository.RefreshTokenRepository) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, passwordHasher: passwordHasher, tokenGen: tokenGen, tokenRepo: tokenRepo}
}

func (u *LoginUseCase) Execute(ctx context.Context, email, password string) (*entity.User, string, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if err := u.passwordHasher.Compare(user.PasswordHash, password); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, err := u.tokenGen.GenerateAccessToken(user.ID, user.Role)
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
