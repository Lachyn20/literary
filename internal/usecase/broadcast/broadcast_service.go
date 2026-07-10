package broadcast

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type BroadcastService struct {
	repo      repository.BroadcastRepository
	fileStore repository.FileStorage
}

func NewBroadcastService(repo repository.BroadcastRepository, fileStore repository.FileStorage) *BroadcastService {
	return &BroadcastService{repo: repo, fileStore: fileStore}
}

func (s *BroadcastService) Create(ctx context.Context, broadcast *entity.Broadcast) error {
	if broadcast.ID == uuid.Nil {
		broadcast.ID = uuid.New()
	}
	return s.repo.Create(ctx, broadcast)
}

func (s *BroadcastService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Broadcast, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BroadcastService) List(ctx context.Context, limit, offset int) ([]*entity.Broadcast, int, error) {
	broadcasts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return broadcasts, len(broadcasts), nil
	}
	
	return broadcasts, total, nil
}

func (s *BroadcastService) Update(ctx context.Context, broadcast *entity.Broadcast) error {
	oldBroadcast, err := s.repo.GetByID(ctx, broadcast.ID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, broadcast); err != nil {
		return err
	}

	if oldBroadcast.FilePath != nil && (broadcast.FilePath == nil || *oldBroadcast.FilePath != *broadcast.FilePath) {
		_ = s.fileStore.Remove(*oldBroadcast.FilePath)
	}

	return nil
}

func (s *BroadcastService) Delete(ctx context.Context, id uuid.UUID) error {
	broadcast, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if broadcast.FilePath != nil && *broadcast.FilePath != "" {
		_ = s.fileStore.Remove(*broadcast.FilePath)
	}

	return nil
}
