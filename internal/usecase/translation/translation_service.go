package translation

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type TranslationService struct {
	byAuthorRepo     repository.TranslatedByAuthorRepository
	intoLanguageRepo repository.TranslatedIntoLanguageRepository
}

func NewTranslationService(byAuthorRepo repository.TranslatedByAuthorRepository, intoLanguageRepo repository.TranslatedIntoLanguageRepository) *TranslationService {
	return &TranslationService{byAuthorRepo: byAuthorRepo, intoLanguageRepo: intoLanguageRepo}
}

// TranslatedByAuthor methods

func (s *TranslationService) CreateByAuthor(ctx context.Context, translation *entity.TranslatedByAuthor) error {
	if translation.ID == uuid.Nil {
		translation.ID = uuid.New()
	}
	return s.byAuthorRepo.Create(ctx, translation)
}

func (s *TranslationService) GetByAuthor(ctx context.Context, id uuid.UUID) (*entity.TranslatedByAuthor, error) {
	return s.byAuthorRepo.GetByID(ctx, id)
}

func (s *TranslationService) ListByAuthor(ctx context.Context) ([]*entity.TranslatedByAuthor, int, error) {
	items, err := s.byAuthorRepo.List(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, len(items), nil
}

func (s *TranslationService) DeleteByAuthor(ctx context.Context, id uuid.UUID) error {
	return s.byAuthorRepo.Delete(ctx, id)
}

// TranslatedIntoLanguage methods

func (s *TranslationService) CreateIntoLanguage(ctx context.Context, translation *entity.TranslatedIntoLanguage) error {
	if translation.ID == uuid.Nil {
		translation.ID = uuid.New()
	}
	return s.intoLanguageRepo.Create(ctx, translation)
}

func (s *TranslationService) GetIntoLanguage(ctx context.Context, id uuid.UUID) (*entity.TranslatedIntoLanguage, error) {
	return s.intoLanguageRepo.GetByID(ctx, id)
}

func (s *TranslationService) ListIntoLanguage(ctx context.Context) ([]*entity.TranslatedIntoLanguage, int, error) {
	items, err := s.intoLanguageRepo.List(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, len(items), nil
}

func (s *TranslationService) DeleteIntoLanguage(ctx context.Context, id uuid.UUID) error {
	return s.intoLanguageRepo.Delete(ctx, id)
}
