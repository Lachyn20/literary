package theatre

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type TheatreService struct {
	repo repository.TheatreProductionRepository
}

func NewTheatreService(repo repository.TheatreProductionRepository) *TheatreService {
	return &TheatreService{repo: repo}
}

func (s *TheatreService) Create(ctx context.Context, theatre *entity.TheatreProduction) error {
	if theatre.ID == uuid.Nil {
		theatre.ID = uuid.New()
	}
	return s.repo.Create(ctx, theatre)
}

func (s *TheatreService) GetByID(ctx context.Context, id uuid.UUID) (*entity.TheatreProduction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TheatreService) List(ctx context.Context) ([]*entity.TheatreProduction, error) {
	return s.repo.List(ctx)
}

func (s *TheatreService) Update(ctx context.Context, theatre *entity.TheatreProduction) error {
	return s.repo.Update(ctx, theatre)
}

func (s *TheatreService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
