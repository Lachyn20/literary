package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, category *entity.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return s.repo.Create(ctx, category)
}

func (s *CategoryService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, limit, offset int) ([]*entity.Category, int, error) {
	categories, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return categories, len(categories), nil
	}
	
	return categories, total, nil
}

func (s *CategoryService) Update(ctx context.Context, category *entity.Category) error {
	return s.repo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
