package work

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateWorkUseCase struct {
	repo repository.WorkRepository
}

func NewCreateWorkUseCase(repo repository.WorkRepository) *CreateWorkUseCase {
	return &CreateWorkUseCase{repo: repo}
}

func (u *CreateWorkUseCase) Execute(ctx context.Context, work *entity.Work) error {
	if work.ID == uuid.Nil {
		work.ID = uuid.New()
	}
	return u.repo.Create(ctx, work)
}

type GetWorkUseCase struct {
	repo repository.WorkRepository
}

func NewGetWorkUseCase(repo repository.WorkRepository) *GetWorkUseCase {
	return &GetWorkUseCase{repo: repo}
}

func (u *GetWorkUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	return u.repo.GetByID(ctx, id)
}

type ListWorksUseCase struct {
	repo repository.WorkRepository
}

func NewListWorksUseCase(repo repository.WorkRepository) *ListWorksUseCase {
	return &ListWorksUseCase{repo: repo}
}

// Execute returns works, total count and error. If filter.Search is set, performs full-text search.
func (u *ListWorksUseCase) Execute(ctx context.Context, filter repository.WorkFilter) ([]*entity.Work, int, error) {
	if filter.Search != nil && *filter.Search != "" {
		return u.repo.Search(ctx, filter)
	}
	works, err := u.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return works, len(works), nil
}

type UpdateWorkUseCase struct {
	repo repository.WorkRepository
}

func NewUpdateWorkUseCase(repo repository.WorkRepository) *UpdateWorkUseCase {
	return &UpdateWorkUseCase{repo: repo}
}

func (u *UpdateWorkUseCase) Execute(ctx context.Context, work *entity.Work) error {
	return u.repo.Update(ctx, work)
}

type DeleteWorkUseCase struct {
	repo repository.WorkRepository
}

func NewDeleteWorkUseCase(repo repository.WorkRepository) *DeleteWorkUseCase {
	return &DeleteWorkUseCase{repo: repo}
}

func (u *DeleteWorkUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
