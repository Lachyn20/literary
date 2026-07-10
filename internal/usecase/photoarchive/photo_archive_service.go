package photoarchive

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type PhotoArchiveService struct {
	repo      repository.PhotoArchiveRepository
	fileStore repository.FileStorage
}

func NewPhotoArchiveService(repo repository.PhotoArchiveRepository, fileStore repository.FileStorage) *PhotoArchiveService {
	return &PhotoArchiveService{repo: repo, fileStore: fileStore}
}

func (s *PhotoArchiveService) Create(ctx context.Context, photo *entity.PhotoArchive) error {
	if photo.ID == uuid.Nil {
		photo.ID = uuid.New()
	}
	return s.repo.Create(ctx, photo)
}

func (s *PhotoArchiveService) GetByID(ctx context.Context, id uuid.UUID) (*entity.PhotoArchive, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PhotoArchiveService) List(ctx context.Context, limit, offset int) ([]*entity.PhotoArchive, int, error) {
	photos, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return photos, len(photos), nil
	}
	
	return photos, total, nil
}

func (s *PhotoArchiveService) Update(ctx context.Context, photo *entity.PhotoArchive) error {
	oldPhoto, err := s.repo.GetByID(ctx, photo.ID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, photo); err != nil {
		return err
	}

	if oldPhoto.ImagePath != "" && oldPhoto.ImagePath != photo.ImagePath {
		_ = s.fileStore.Remove(oldPhoto.ImagePath)
	}

	return nil
}

func (s *PhotoArchiveService) Delete(ctx context.Context, id uuid.UUID) error {
	photo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if photo.ImagePath != "" {
		_ = s.fileStore.Remove(photo.ImagePath)
	}

	return nil
}
