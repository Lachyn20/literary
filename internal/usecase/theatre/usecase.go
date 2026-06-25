package theatre

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateTheatreProductionUseCase struct {
	repo repository.TheatreProductionRepository
}

func NewCreateTheatreProductionUseCase(repo repository.TheatreProductionRepository) *CreateTheatreProductionUseCase {
	return &CreateTheatreProductionUseCase{repo: repo}
}

func (u *CreateTheatreProductionUseCase) Execute(ctx context.Context, production *entity.TheatreProduction) error {
	if production.ID == uuid.Nil {
		production.ID = uuid.New()
	}
	return u.repo.Create(ctx, production)
}

type GetTheatreProductionUseCase struct {
	repo repository.TheatreProductionRepository
}

func NewGetTheatreProductionUseCase(repo repository.TheatreProductionRepository) *GetTheatreProductionUseCase {
	return &GetTheatreProductionUseCase{repo: repo}
}

func (u *GetTheatreProductionUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.TheatreProduction, error) {
	return u.repo.GetByID(ctx, id)
}

type ListTheatreProductionsUseCase struct {
	repo repository.TheatreProductionRepository
}

func NewListTheatreProductionsUseCase(repo repository.TheatreProductionRepository) *ListTheatreProductionsUseCase {
	return &ListTheatreProductionsUseCase{repo: repo}
}

func (u *ListTheatreProductionsUseCase) Execute(ctx context.Context) ([]*entity.TheatreProduction, error) {
	return u.repo.List(ctx)
}

type UpdateTheatreProductionUseCase struct {
	repo repository.TheatreProductionRepository
}

func NewUpdateTheatreProductionUseCase(repo repository.TheatreProductionRepository) *UpdateTheatreProductionUseCase {
	return &UpdateTheatreProductionUseCase{repo: repo}
}

func (u *UpdateTheatreProductionUseCase) Execute(ctx context.Context, production *entity.TheatreProduction) error {
	return u.repo.Update(ctx, production)
}

type DeleteTheatreProductionUseCase struct {
	repo repository.TheatreProductionRepository
}

func NewDeleteTheatreProductionUseCase(repo repository.TheatreProductionRepository) *DeleteTheatreProductionUseCase {
	return &DeleteTheatreProductionUseCase{repo: repo}
}

func (u *DeleteTheatreProductionUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
