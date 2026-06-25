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

type BookPhotoRepository struct {
	pool *pgxpool.Pool
}

func NewBookPhotoRepository(pool *pgxpool.Pool) *BookPhotoRepository {
	return &BookPhotoRepository{pool: pool}
}

func (r *BookPhotoRepository) Create(ctx context.Context, photo *entity.BookPhoto) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO book_photos (id, book_id, image_path) VALUES ($1,$2,$3)`, photo.ID, photo.BookID, photo.ImagePath)
	return err
}

func (r *BookPhotoRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.BookPhoto, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,book_id,image_path FROM book_photos WHERE id=$1`, id)
	var photo entity.BookPhoto
	if err := row.Scan(&photo.ID, &photo.BookID, &photo.ImagePath); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &photo, nil
}

func (r *BookPhotoRepository) ListByBookID(ctx context.Context, bookID uuid.UUID) ([]*entity.BookPhoto, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,book_id,image_path FROM book_photos WHERE book_id=$1`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*entity.BookPhoto
	for rows.Next() {
		var photo entity.BookPhoto
		if err := rows.Scan(&photo.ID, &photo.BookID, &photo.ImagePath); err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}
	return photos, rows.Err()
}

func (r *BookPhotoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM book_photos WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
