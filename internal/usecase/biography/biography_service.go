package biography

import (
	"context"
	"errors"
	"time"

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

type BiographyEventService struct {
	repo    repository.BiographyEventRepository
	bioRepo repository.BiographyRepository
}

func NewBiographyEventService(repo repository.BiographyEventRepository, bioRepo repository.BiographyRepository) *BiographyEventService {
	return &BiographyEventService{repo: repo, bioRepo: bioRepo}
}

func (s *BiographyEventService) Create(ctx context.Context, event *entity.BiographyEvent) error {
	// Ensure biography exists
	bio, err := s.bioRepo.GetLatest(ctx)
	if err != nil {
		return errors.New("biography not found, create biography first")
	}
	if bio == nil {
		return errors.New("biography not found, create biography first")
	}

	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	event.BiographyID = bio.ID
	event.CreatedAt = time.Now()
	return s.repo.Create(ctx, event)
}

func (s *BiographyEventService) GetByID(ctx context.Context, id uuid.UUID) (*entity.BiographyEvent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BiographyEventService) List(ctx context.Context, biographyID uuid.UUID) ([]*entity.BiographyEvent, error) {
	return s.repo.ListByBiographyID(ctx, biographyID)
}

func (s *BiographyEventService) Update(ctx context.Context, event *entity.BiographyEvent) error {
	return s.repo.Update(ctx, event)
}

func (s *BiographyEventService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
