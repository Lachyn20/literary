package auth

import (
	"context"

	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type LogoutUseCase struct {
	tokenRepo repository.RefreshTokenRepository
}

func NewLogoutUseCase(tokenRepo repository.RefreshTokenRepository) *LogoutUseCase {
	return &LogoutUseCase{tokenRepo: tokenRepo}
}

func (u *LogoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	rt, err := u.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return u.tokenRepo.Delete(ctx, rt.ID)
}
