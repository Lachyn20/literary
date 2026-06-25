package category

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateCategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewCreateCategoryUseCase(repo repository.CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo}
}

func (u *CreateCategoryUseCase) Execute(ctx context.Context, category *entity.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	if category.Slug == "" {
		category.Slug = slugifyCategoryName(category.Name)
	}
	return u.repo.Create(ctx, category)
}

type GetCategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewGetCategoryUseCase(repo repository.CategoryRepository) *GetCategoryUseCase {
	return &GetCategoryUseCase{repo: repo}
}

func (u *GetCategoryUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return u.repo.GetByID(ctx, id)
}

type ListCategoriesUseCase struct {
	repo repository.CategoryRepository
}

func NewListCategoriesUseCase(repo repository.CategoryRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{repo: repo}
}

func (u *ListCategoriesUseCase) Execute(ctx context.Context) ([]*entity.Category, error) {
	return u.repo.List(ctx)
}

type UpdateCategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewUpdateCategoryUseCase(repo repository.CategoryRepository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{repo: repo}
}

func (u *UpdateCategoryUseCase) Execute(ctx context.Context, category *entity.Category) error {
	if category.Slug == "" {
		category.Slug = slugifyCategoryName(category.Name)
	}
	return u.repo.Update(ctx, category)
}

type DeleteCategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewDeleteCategoryUseCase(repo repository.CategoryRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{repo: repo}
}

func (u *DeleteCategoryUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func slugifyCategoryName(name string) string {
	slug := strings.TrimSpace(strings.ToLower(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, "--", "-")
	return slug
}
