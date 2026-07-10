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

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO categories (id, name, slug) VALUES ($1,$2,$3)`, category.ID, category.Name, category.Slug)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,name,slug FROM categories WHERE id=$1`, id)
	var category entity.Category
	if err := row.Scan(&category.ID, &category.Name, &category.Slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,name,slug FROM categories WHERE slug=$1`, slug)
	var category entity.Category
	if err := row.Scan(&category.ID, &category.Name, &category.Slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) List(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,name,slug FROM categories ORDER BY name ASC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM categories`).Scan(&count)
	return count, err
}

func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE categories SET name=$1, slug=$2 WHERE id=$3`, category.Name, category.Slug, category.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
