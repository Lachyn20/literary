package film

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateFilmUseCase struct {
	repo repository.FilmRepository
}

func NewCreateFilmUseCase(repo repository.FilmRepository) *CreateFilmUseCase {
	return &CreateFilmUseCase{repo: repo}
}

func (u *CreateFilmUseCase) Execute(ctx context.Context, film *entity.Film) error {
	if film.ID == uuid.Nil {
		film.ID = uuid.New()
	}
	return u.repo.Create(ctx, film)
}

type GetFilmUseCase struct {
	repo repository.FilmRepository
}

func NewGetFilmUseCase(repo repository.FilmRepository) *GetFilmUseCase {
	return &GetFilmUseCase{repo: repo}
}

func (u *GetFilmUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Film, error) {
	return u.repo.GetByID(ctx, id)
}

type ListFilmsUseCase struct {
	repo repository.FilmRepository
}

func NewListFilmsUseCase(repo repository.FilmRepository) *ListFilmsUseCase {
	return &ListFilmsUseCase{repo: repo}
}

func (u *ListFilmsUseCase) Execute(ctx context.Context) ([]*entity.Film, error) {
	return u.repo.List(ctx)
}

type UpdateFilmUseCase struct {
	repo repository.FilmRepository
}

func NewUpdateFilmUseCase(repo repository.FilmRepository) *UpdateFilmUseCase {
	return &UpdateFilmUseCase{repo: repo}
}

func (u *UpdateFilmUseCase) Execute(ctx context.Context, film *entity.Film) error {
	return u.repo.Update(ctx, film)
}

type DeleteFilmUseCase struct {
	repo repository.FilmRepository
}

func NewDeleteFilmUseCase(repo repository.FilmRepository) *DeleteFilmUseCase {
	return &DeleteFilmUseCase{repo: repo}
}

func (u *DeleteFilmUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
