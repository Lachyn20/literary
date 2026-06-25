package biography

import (
	"context"

	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type GetBiographyUseCase struct {
	repo repository.BiographyRepository
}

func NewGetBiographyUseCase(repo repository.BiographyRepository) *GetBiographyUseCase {
	return &GetBiographyUseCase{repo: repo}
}

func (u *GetBiographyUseCase) Execute(ctx context.Context) (*entity.Biography, error) {
	return u.repo.GetLatest(ctx)
}

type UpdateBiographyUseCase struct {
	repo repository.BiographyRepository
}

func NewUpdateBiographyUseCase(repo repository.BiographyRepository) *UpdateBiographyUseCase {
	return &UpdateBiographyUseCase{repo: repo}
}

func (u *UpdateBiographyUseCase) Execute(ctx context.Context, biography *entity.Biography) error {
	return u.repo.Update(ctx, biography)
}
type CreateBiographyUseCase struct {
	repo repository.BiographyRepository
}

func NewCreateBiographyUseCase(repo repository.BiographyRepository) *CreateBiographyUseCase {
	return &CreateBiographyUseCase{repo: repo}
}

func (u *CreateBiographyUseCase) Execute(ctx context.Context, biography *entity.Biography) error {
	return u.repo.Create(ctx, biography)
}