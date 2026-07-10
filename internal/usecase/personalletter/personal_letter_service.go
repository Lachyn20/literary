package personalletter

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type PersonalLetterService struct {
	repo      repository.PersonalLetterRepository
	fileStore repository.FileStorage
}

func NewPersonalLetterService(repo repository.PersonalLetterRepository, fileStore repository.FileStorage) *PersonalLetterService {
	return &PersonalLetterService{repo: repo, fileStore: fileStore}
}

func (s *PersonalLetterService) Create(ctx context.Context, letter *entity.PersonalLetter) error {
	if letter.ID == uuid.Nil {
		letter.ID = uuid.New()
	}
	return s.repo.Create(ctx, letter)
}

func (s *PersonalLetterService) GetByID(ctx context.Context, id uuid.UUID) (*entity.PersonalLetter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PersonalLetterService) List(ctx context.Context, limit, offset int) ([]*entity.PersonalLetter, int, error) {
	letters, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return letters, len(letters), nil
	}
	
	return letters, total, nil
}

func (s *PersonalLetterService) Update(ctx context.Context, letter *entity.PersonalLetter) error {
	oldLetter, err := s.repo.GetByID(ctx, letter.ID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, letter); err != nil {
		return err
	}

	if oldLetter.ScanImagePath != nil && (letter.ScanImagePath == nil || *oldLetter.ScanImagePath != *letter.ScanImagePath) {
		_ = s.fileStore.Remove(*oldLetter.ScanImagePath)
	}

	return nil
}

func (s *PersonalLetterService) Delete(ctx context.Context, id uuid.UUID) error {
	letter, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if letter.ScanImagePath != nil && *letter.ScanImagePath != "" {
		_ = s.fileStore.Remove(*letter.ScanImagePath)
	}

	return nil
}
