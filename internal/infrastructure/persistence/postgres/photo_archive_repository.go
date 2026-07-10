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

type PhotoArchiveRepository struct {
	pool *pgxpool.Pool
}

func NewPhotoArchiveRepository(pool *pgxpool.Pool) *PhotoArchiveRepository {
	return &PhotoArchiveRepository{pool: pool}
}

func (r *PhotoArchiveRepository) Create(ctx context.Context, photo *entity.PhotoArchive) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO photo_archive (id, title, image_path, description, taken_date, category, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`, photo.ID, photo.Title, photo.ImagePath, photo.Description, photo.TakenDate, photo.Category, photo.CreatedAt)
	return err
}

func (r *PhotoArchiveRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PhotoArchive, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,image_path,description,taken_date,category,created_at FROM photo_archive WHERE id=$1`, id)
	var photo entity.PhotoArchive
	if err := row.Scan(&photo.ID, &photo.Title, &photo.ImagePath, &photo.Description, &photo.TakenDate, &photo.Category, &photo.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &photo, nil
}

func (r *PhotoArchiveRepository) List(ctx context.Context, limit, offset int) ([]*entity.PhotoArchive, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,image_path,description,taken_date,category,created_at FROM photo_archive ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*entity.PhotoArchive
	for rows.Next() {
		var photo entity.PhotoArchive
		if err := rows.Scan(&photo.ID, &photo.Title, &photo.ImagePath, &photo.Description, &photo.TakenDate, &photo.Category, &photo.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}
	return photos, rows.Err()
}

func (r *PhotoArchiveRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM photo_archive`).Scan(&count)
	return count, err
}

func (r *PhotoArchiveRepository) Update(ctx context.Context, photo *entity.PhotoArchive) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE photo_archive SET title=$1,image_path=$2,description=$3,taken_date=$4,category=$5 WHERE id=$6`, photo.Title, photo.ImagePath, photo.Description, photo.TakenDate, photo.Category, photo.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *PhotoArchiveRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM photo_archive WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
