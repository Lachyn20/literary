package work

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type WorkService struct {
	repo      repository.WorkRepository
	fileStore repository.FileStorage
}

func NewWorkService(repo repository.WorkRepository, fileStore repository.FileStorage) *WorkService {
	return &WorkService{repo: repo, fileStore: fileStore}
}

func (s *WorkService) Create(ctx context.Context, work *entity.Work) error {
	if work.ID == uuid.Nil {
		work.ID = uuid.New()
	}
	return s.repo.Create(ctx, work)
}

func (s *WorkService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WorkService) List(ctx context.Context, filter repository.WorkFilter) ([]*entity.Work, int, error) {
	if filter.Search != nil && *filter.Search != "" {
		return s.repo.Search(ctx, filter)
	}
	works, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return works, len(works), nil
}

func (s *WorkService) Update(ctx context.Context, work *entity.Work) error {
	oldWork, err := s.repo.GetByID(ctx, work.ID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, work); err != nil {
		return err
	}

	if oldWork.FilePath != nil && (work.FilePath == nil || *oldWork.FilePath != *work.FilePath) {
		_ = s.fileStore.Remove(*oldWork.FilePath)
	}

	return nil
}

func (s *WorkService) Delete(ctx context.Context, id uuid.UUID) error {
	work, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if work.FilePath != nil && *work.FilePath != "" {
		_ = s.fileStore.Remove(*work.FilePath)
	}

	return nil
}
