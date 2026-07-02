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
	repo      repository.WorkRepository
	fileStore repository.FileStorage
}

func NewUpdateWorkUseCase(repo repository.WorkRepository, fileStore repository.FileStorage) *UpdateWorkUseCase {
	return &UpdateWorkUseCase{repo: repo, fileStore: fileStore}
}

func (u *UpdateWorkUseCase) Execute(ctx context.Context, work *entity.Work) error {
	oldWork, err := u.repo.GetByID(ctx, work.ID)
	if err != nil {
		return err
	}

	if err := u.repo.Update(ctx, work); err != nil {
		return err
	}

	if oldWork.FilePath != nil && (work.FilePath == nil || *oldWork.FilePath != *work.FilePath) {
		_ = u.fileStore.Remove(*oldWork.FilePath)
	}

	return nil
}

type DeleteWorkUseCase struct {
	repo      repository.WorkRepository
	fileStore repository.FileStorage
}

func NewDeleteWorkUseCase(repo repository.WorkRepository, fileStore repository.FileStorage) *DeleteWorkUseCase {
	return &DeleteWorkUseCase{repo: repo, fileStore: fileStore}
}

func (u *DeleteWorkUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	work, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	if work.FilePath != nil && *work.FilePath != "" {
		_ = u.fileStore.Remove(*work.FilePath)
	}

	return nil
}
