package postgres

import (
	"context"
	"errors"

	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BiographyRepository struct {
	pool *pgxpool.Pool
}

func NewBiographyRepository(pool *pgxpool.Pool) *BiographyRepository {
	return &BiographyRepository{pool: pool}
}

func (r *BiographyRepository) Create(ctx context.Context, biography *entity.Biography) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO biography (id, photo_path, updated_at) VALUES ($1,$2,$3)`, biography.ID, biography.PhotoPath, biography.UpdatedAt)
	return err
}

func (r *BiographyRepository) GetLatest(ctx context.Context) (*entity.Biography, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, photo_path, updated_at FROM biography ORDER BY updated_at DESC LIMIT 1`)
	var biography entity.Biography
	if err := row.Scan(&biography.ID, &biography.PhotoPath, &biography.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// Load events for the biography
	eventRepo := NewBiographyEventRepository(r.pool)
	events, err := eventRepo.ListByBiographyID(ctx, biography.ID)
	if err == nil {
		biography.Events = make([]entity.BiographyEvent, 0, len(events))
		for _, e := range events {
			biography.Events = append(biography.Events, *e)
		}
	}

	return &biography, nil
}

func (r *BiographyRepository) Update(ctx context.Context, biography *entity.Biography) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE biography SET photo_path=$1, updated_at=$2 WHERE id=$3`, biography.PhotoPath, biography.UpdatedAt, biography.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
