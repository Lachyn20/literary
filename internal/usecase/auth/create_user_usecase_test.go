package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type stubUserRepo struct {
	users map[string]*entity.User
}

func newStubUserRepo() *stubUserRepo {
	return &stubUserRepo{users: make(map[string]*entity.User)}
}

func (s *stubUserRepo) Create(ctx context.Context, user *entity.User) error {
	s.users[user.Email] = user
	return nil
}

func (s *stubUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (s *stubUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, ok := s.users[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (s *stubUserRepo) List(ctx context.Context) ([]*entity.User, error) {
	out := make([]*entity.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out, nil
}

func (s *stubUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	for email, u := range s.users {
		if u.ID == id {
			delete(s.users, email)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (s *stubUserRepo) Update(ctx context.Context, user *entity.User) error {
	for email, u := range s.users {
		if u.ID == user.ID {
			s.users[email] = user
			return nil
		}
	}
	return repository.ErrNotFound
}

func (s *stubUserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	for _, u := range s.users {
		if u.ID == id {
			u.PasswordHash = passwordHash
			return nil
		}
	}
	return repository.ErrNotFound
}

type stubHasher struct{}

func (stubHasher) Hash(password string) (string, error) { return "hashed:" + password, nil }
func (stubHasher) Compare(hash, password string) error {
	if hash == "hashed:"+password {
		return nil
	}
	return errors.New("mismatch")
}

func TestCreateUserUseCaseRejectsNonAdminCaller(t *testing.T) {
	repo := newStubUserRepo()
	usecase := NewCreateUserUseCase(repo, stubHasher{}, nil, nil)

	_, err := usecase.Execute(context.Background(), "editor", "Jane", "jane@example.com", "secret123", entity.RoleEditor)
	if err == nil {
		t.Fatal("expected non-admin caller to be rejected")
	}
	if err.Error() != "only admins can create users" {
		t.Fatalf("unexpected error: %v", err)
	}
}
