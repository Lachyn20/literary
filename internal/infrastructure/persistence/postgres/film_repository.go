package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type FilmRepository struct {
	pool *pgxpool.Pool
}

func NewFilmRepository(pool *pgxpool.Pool) *FilmRepository {
	return &FilmRepository{pool: pool}
}

func (r *FilmRepository) Create(ctx context.Context, film *entity.Film) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO films (id, title, film_type, based_on_scenario, director, release_year, video_path, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, film.ID, film.Title, film.FilmType, film.BasedOnScenario, film.Director, film.ReleaseYear, film.VideoPath, film.CreatedAt)
	return err
}

func (r *FilmRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Film, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,film_type,based_on_scenario,director,release_year,video_path,created_at FROM films WHERE id=$1`, id)
	var film entity.Film
	if err := row.Scan(&film.ID, &film.Title, &film.FilmType, &film.BasedOnScenario, &film.Director, &film.ReleaseYear, &film.VideoPath, &film.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &film, nil
}

func (r *FilmRepository) List(ctx context.Context) ([]*entity.Film, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,film_type,based_on_scenario,director,release_year,video_path,created_at FROM films`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []*entity.Film
	for rows.Next() {
		var film entity.Film
		if err := rows.Scan(&film.ID, &film.Title, &film.FilmType, &film.BasedOnScenario, &film.Director, &film.ReleaseYear, &film.VideoPath, &film.CreatedAt); err != nil {
			return nil, err
		}
		films = append(films, &film)
	}
	return films, rows.Err()
}

func (r *FilmRepository) Update(ctx context.Context, film *entity.Film) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE films SET title=$1, film_type=$2, based_on_scenario=$3, director=$4, release_year=$5, video_path=$6 WHERE id=$7`, film.Title, film.FilmType, film.BasedOnScenario, film.Director, film.ReleaseYear, film.VideoPath, film.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *FilmRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM films WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
