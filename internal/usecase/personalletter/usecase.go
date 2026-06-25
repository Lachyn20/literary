package personalletter

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreatePersonalLetterUseCase struct {
	repo repository.PersonalLetterRepository
}

func NewCreatePersonalLetterUseCase(repo repository.PersonalLetterRepository) *CreatePersonalLetterUseCase {
	return &CreatePersonalLetterUseCase{repo: repo}
}

func (u *CreatePersonalLetterUseCase) Execute(ctx context.Context, letter *entity.PersonalLetter) error {
	if letter.ID == uuid.Nil {
		letter.ID = uuid.New()
	}
	return u.repo.Create(ctx, letter)
}

type GetPersonalLetterUseCase struct {
	repo repository.PersonalLetterRepository
}

func NewGetPersonalLetterUseCase(repo repository.PersonalLetterRepository) *GetPersonalLetterUseCase {
	return &GetPersonalLetterUseCase{repo: repo}
}

func (u *GetPersonalLetterUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.PersonalLetter, error) {
	return u.repo.GetByID(ctx, id)
}

type ListPersonalLettersUseCase struct {
	repo repository.PersonalLetterRepository
}

func NewListPersonalLettersUseCase(repo repository.PersonalLetterRepository) *ListPersonalLettersUseCase {
	return &ListPersonalLettersUseCase{repo: repo}
}

func (u *ListPersonalLettersUseCase) Execute(ctx context.Context) ([]*entity.PersonalLetter, error) {
	return u.repo.List(ctx)
}

type UpdatePersonalLetterUseCase struct {
	repo repository.PersonalLetterRepository
}

func NewUpdatePersonalLetterUseCase(repo repository.PersonalLetterRepository) *UpdatePersonalLetterUseCase {
	return &UpdatePersonalLetterUseCase{repo: repo}
}

func (u *UpdatePersonalLetterUseCase) Execute(ctx context.Context, letter *entity.PersonalLetter) error {
	return u.repo.Update(ctx, letter)
}

type DeletePersonalLetterUseCase struct {
	repo repository.PersonalLetterRepository
}

func NewDeletePersonalLetterUseCase(repo repository.PersonalLetterRepository) *DeletePersonalLetterUseCase {
	return &DeletePersonalLetterUseCase{repo: repo}
}

func (u *DeletePersonalLetterUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
