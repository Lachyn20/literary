package work

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type mockWorkRepo struct{
	created bool
}

func (m *mockWorkRepo) Create(ctx context.Context, work *entity.Work) error {
	m.created = true
	if work.ID == uuid.Nil {
		work.ID = uuid.New()
	}
	return nil
}
func (m *mockWorkRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) { return nil, repository.ErrNotFound }
func (m *mockWorkRepo) List(ctx context.Context, filter repository.WorkFilter) ([]*entity.Work, error) { return nil, nil }
func (m *mockWorkRepo) Search(ctx context.Context, filter repository.WorkFilter) ([]*entity.Work, int, error) { return nil, 0, nil }
func (m *mockWorkRepo) Update(ctx context.Context, work *entity.Work) error { return nil }
func (m *mockWorkRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

func TestCreateWorkUseCase(t *testing.T) {
	m := &mockWorkRepo{}
	u := NewCreateWorkUseCase(m)
	w := &entity.Work{Title: "Test", AudienceType: entity.AudienceType("adult")}
	if err := u.Execute(context.Background(), w); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !m.created {
		t.Fatalf("expected repo Create to be called")
	}
	if w.ID == uuid.Nil {
		t.Fatalf("expected ID to be set")
	}
}
