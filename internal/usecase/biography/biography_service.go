package biography

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type BiographyService struct {
	repo      repository.BiographyRepository
	fileStore repository.FileStorage
}

func NewBiographyService(repo repository.BiographyRepository, fileStore repository.FileStorage) *BiographyService {
	return &BiographyService{repo: repo, fileStore: fileStore}
}

func (s *BiographyService) Create(ctx context.Context, biography *entity.Biography) error {
	// Check if biography already exists
	existing, err := s.repo.GetLatest(ctx)
	if err == nil && existing != nil {
		return errors.New("biography already exists, use update instead")
	}

	if biography.ID == uuid.Nil {
		biography.ID = uuid.New()
	}
	return s.repo.Create(ctx, biography)
}

func (s *BiographyService) GetLatest(ctx context.Context) (*entity.Biography, error) {
	return s.repo.GetLatest(ctx)
}

func (s *BiographyService) Update(ctx context.Context, biography *entity.Biography) error {
	old, err := s.repo.GetLatest(ctx)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, biography); err != nil {
		return err
	}

	if old.PhotoPath != nil && (biography.PhotoPath == nil || *old.PhotoPath != *biography.PhotoPath) {
		_ = s.fileStore.Remove(*old.PhotoPath)
	}

	return nil
}
