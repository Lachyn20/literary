package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TheatreProductionRepository struct {
	pool *pgxpool.Pool
}

func NewTheatreProductionRepository(pool *pgxpool.Pool) *TheatreProductionRepository {
	return &TheatreProductionRepository{pool: pool}
}

func (r *TheatreProductionRepository) Create(ctx context.Context, production *entity.TheatreProduction) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO theatre_productions (id, play_title, theatre_name, premiere_date, notes, created_at) VALUES ($1,$2,$3,$4,$5,$6)`, production.ID, production.PlayTitle, production.TheatreName, production.PremiereDate, production.Notes, production.CreatedAt)
	return err
}

func (r *TheatreProductionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TheatreProduction, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,play_title,theatre_name,premiere_date,notes,created_at FROM theatre_productions WHERE id=$1`, id)
	var production entity.TheatreProduction
	if err := row.Scan(&production.ID, &production.PlayTitle, &production.TheatreName, &production.PremiereDate, &production.Notes, &production.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &production, nil
}

func (r *TheatreProductionRepository) List(ctx context.Context, limit, offset int) ([]*entity.TheatreProduction, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,play_title,theatre_name,premiere_date,notes,created_at FROM theatre_productions ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productions []*entity.TheatreProduction
	for rows.Next() {
		var production entity.TheatreProduction
		if err := rows.Scan(&production.ID, &production.PlayTitle, &production.TheatreName, &production.PremiereDate, &production.Notes, &production.CreatedAt); err != nil {
			return nil, err
		}
		productions = append(productions, &production)
	}
	return productions, rows.Err()
}

func (r *TheatreProductionRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM theatre_productions`).Scan(&count)
	return count, err
}

func (r *TheatreProductionRepository) Update(ctx context.Context, production *entity.TheatreProduction) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE theatre_productions SET play_title=$1,theatre_name=$2,premiere_date=$3,notes=$4 WHERE id=$5`, production.PlayTitle, production.TheatreName, production.PremiereDate, production.Notes, production.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *TheatreProductionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM theatre_productions WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
