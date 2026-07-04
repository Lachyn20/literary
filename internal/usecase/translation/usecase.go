package translation

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateTranslatedByAuthorUseCase struct {
	repo repository.TranslatedByAuthorRepository
}

func NewCreateTranslatedByAuthorUseCase(repo repository.TranslatedByAuthorRepository) *CreateTranslatedByAuthorUseCase {
	return &CreateTranslatedByAuthorUseCase{repo: repo}
}

func (u *CreateTranslatedByAuthorUseCase) Execute(ctx context.Context, translation *entity.TranslatedByAuthor) error {
	if translation.ID == uuid.Nil {
		translation.ID = uuid.New()
	}
	return u.repo.Create(ctx, translation)
}

type GetTranslatedByAuthorUseCase struct {
	repo repository.TranslatedByAuthorRepository
}

func NewGetTranslatedByAuthorUseCase(repo repository.TranslatedByAuthorRepository) *GetTranslatedByAuthorUseCase {
	return &GetTranslatedByAuthorUseCase{repo: repo}
}

func (u *GetTranslatedByAuthorUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.TranslatedByAuthor, error) {
	return u.repo.GetByID(ctx, id)
}

type ListTranslatedByAuthorsUseCase struct {
	repo repository.TranslatedByAuthorRepository
}

func NewListTranslatedByAuthorsUseCase(repo repository.TranslatedByAuthorRepository) *ListTranslatedByAuthorsUseCase {
	return &ListTranslatedByAuthorsUseCase{repo: repo}
}

func (u *ListTranslatedByAuthorsUseCase) Execute(ctx context.Context) ([]*entity.TranslatedByAuthor, error) {
	return u.repo.List(ctx)
}

	type DeleteTranslatedByAuthorUseCase struct {
	repo repository.TranslatedByAuthorRepository
}

func NewDeleteTranslatedByAuthorUseCase(repo repository.TranslatedByAuthorRepository) *DeleteTranslatedByAuthorUseCase {
	return &DeleteTranslatedByAuthorUseCase{repo: repo}
}

func (u *DeleteTranslatedByAuthorUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

type CreateTranslatedIntoLanguageUseCase struct {
	repo repository.TranslatedIntoLanguageRepository
}

func NewCreateTranslatedIntoLanguageUseCase(repo repository.TranslatedIntoLanguageRepository) *CreateTranslatedIntoLanguageUseCase {
	return &CreateTranslatedIntoLanguageUseCase{repo: repo}
}

func (u *CreateTranslatedIntoLanguageUseCase) Execute(ctx context.Context, translation *entity.TranslatedIntoLanguage) error {
	if translation.ID == uuid.Nil {
		translation.ID = uuid.New()
	}
	return u.repo.Create(ctx, translation)
}

type GetTranslatedIntoLanguageUseCase struct {
	repo repository.TranslatedIntoLanguageRepository
}

func NewGetTranslatedIntoLanguageUseCase(repo repository.TranslatedIntoLanguageRepository) *GetTranslatedIntoLanguageUseCase {
	return &GetTranslatedIntoLanguageUseCase{repo: repo}
}

func (u *GetTranslatedIntoLanguageUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.TranslatedIntoLanguage, error) {
	return u.repo.GetByID(ctx, id)
}

type ListTranslatedIntoLanguagesUseCase struct {
	repo repository.TranslatedIntoLanguageRepository
}

func NewListTranslatedIntoLanguagesUseCase(repo repository.TranslatedIntoLanguageRepository) *ListTranslatedIntoLanguagesUseCase {
	return &ListTranslatedIntoLanguagesUseCase{repo: repo}
}

func (u *ListTranslatedIntoLanguagesUseCase) Execute(ctx context.Context) ([]*entity.TranslatedIntoLanguage, error) {
	return u.repo.List(ctx)
}

type DeleteTranslatedIntoLanguageUseCase struct {
	repo repository.TranslatedIntoLanguageRepository
}

func NewDeleteTranslatedIntoLanguageUseCase(repo repository.TranslatedIntoLanguageRepository) *DeleteTranslatedIntoLanguageUseCase {
	return &DeleteTranslatedIntoLanguageUseCase{repo: repo}
}

func (u *DeleteTranslatedIntoLanguageUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
