package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type RefreshTokenUseCase struct {
	userRepo  repository.UserRepository
	tokenGen  repository.TokenGenerator
	tokenRepo repository.RefreshTokenRepository
}

func NewRefreshTokenUseCase(userRepo repository.UserRepository, tokenGen repository.TokenGenerator, tokenRepo repository.RefreshTokenRepository) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{userRepo: userRepo, tokenGen: tokenGen, tokenRepo: tokenRepo}
}

func (u *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (string, string, error) {
	rt, err := u.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token expired")
	}

	user, err := u.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}

	accessToken, err := u.tokenGen.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := u.tokenGen.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := u.tokenRepo.Delete(ctx, rt.ID); err != nil {
		return "", "", err
	}

	if err := u.tokenRepo.Create(ctx, &entity.RefreshToken{ID: uuid.New(), UserID: user.ID, Token: newRefresh, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), CreatedAt: time.Now()}); err != nil {
		return "", "", err
	}

	return accessToken, newRefresh, nil
}
