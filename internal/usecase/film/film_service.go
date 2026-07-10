package film

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type FilmService struct {
	repo      repository.FilmRepository
	fileStore repository.FileStorage
}

func NewFilmService(repo repository.FilmRepository, fileStore repository.FileStorage) *FilmService {
	return &FilmService{repo: repo, fileStore: fileStore}
}

func (s *FilmService) Create(ctx context.Context, film *entity.Film) error {
	if film.ID == uuid.Nil {
		film.ID = uuid.New()
	}
	return s.repo.Create(ctx, film)
}

func (s *FilmService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Film, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FilmService) List(ctx context.Context, limit, offset int) ([]*entity.Film, int, error) {
	films, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return films, len(films), nil
	}
	
	return films, total, nil
}

func (s *FilmService) Update(ctx context.Context, film *entity.Film) error {
	oldFilm, err := s.repo.GetByID(ctx, film.ID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, film); err != nil {
		return err
	}

	if oldFilm.VideoPath != nil && (film.VideoPath == nil || *oldFilm.VideoPath != *film.VideoPath) {
		_ = s.fileStore.Remove(*oldFilm.VideoPath)
	}

	return nil
}

func (s *FilmService) Delete(ctx context.Context, id uuid.UUID) error {
	film, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if film.VideoPath != nil && *film.VideoPath != "" {
		_ = s.fileStore.Remove(*film.VideoPath)
	}

	return nil
}
