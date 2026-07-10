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
		Active:       true,
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

type CreateUserUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher repository.PasswordHasher
	tokenGen       repository.TokenGenerator
	tokenRepo      repository.RefreshTokenRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository, passwordHasher repository.PasswordHasher, tokenGen repository.TokenGenerator, tokenRepo repository.RefreshTokenRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo, passwordHasher: passwordHasher, tokenGen: tokenGen, tokenRepo: tokenRepo}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, callerRole, name, email, password, role string) (*entity.User, error) {
	if callerRole != entity.RoleAdmin {
		return nil, errors.New("only admins can create users")
	}
	if role != entity.RoleAdmin && role != entity.RoleEditor {
		return nil, errors.New("invalid role")
	}
	if _, err := u.userRepo.GetByEmail(ctx, email); err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{ID: uuid.New(), Name: name, Email: email, PasswordHash: hash, Role: role, Active: true, CreatedAt: time.Now()}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *CreateUserUseCase) ListUsers(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.List(ctx)
}

func (u *CreateUserUseCase) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.Active = active
	return u.userRepo.Update(ctx, user)
}

type ChangePasswordUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher repository.PasswordHasher
}

func NewChangePasswordUseCase(userRepo repository.UserRepository, passwordHasher repository.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userRepo: userRepo, passwordHasher: passwordHasher}
}

func (u *ChangePasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := u.passwordHasher.Compare(user.PasswordHash, currentPassword); err != nil {
		return errors.New("current password is invalid")
	}
	hash, err := u.passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}
	return u.userRepo.UpdatePassword(ctx, userID, hash)
}
