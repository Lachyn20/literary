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

type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

func (r *BookRepository) Create(ctx context.Context, book *entity.Book) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO books (id, title, bibliographic_info, cover_image_path, pdf_path, page_count, published_year, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, book.ID, book.Title, book.BibliographicInfo, book.CoverImagePath, book.PDFPath, book.PageCount, book.PublishedYear, book.CreatedAt)
	return err
}

func (r *BookRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,bibliographic_info,cover_image_path,pdf_path,page_count,published_year,created_at FROM books WHERE id=$1`, id)
	var book entity.Book
	if err := row.Scan(&book.ID, &book.Title, &book.BibliographicInfo, &book.CoverImagePath, &book.PDFPath, &book.PageCount, &book.PublishedYear, &book.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) List(ctx context.Context, limit, offset int) ([]*entity.Book, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,bibliographic_info,cover_image_path,pdf_path,page_count,published_year,created_at FROM books ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entity.Book
	for rows.Next() {
		var book entity.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.BibliographicInfo, &book.CoverImagePath, &book.PDFPath, &book.PageCount, &book.PublishedYear, &book.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	return books, rows.Err()
}

func (r *BookRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM books`).Scan(&count)
	return count, err
}

func (r *BookRepository) Update(ctx context.Context, book *entity.Book) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE books SET title=$1, bibliographic_info=$2, cover_image_path=$3, pdf_path=$4, page_count=$5, published_year=$6 WHERE id=$7`, book.Title, book.BibliographicInfo, book.CoverImagePath, book.PDFPath, book.PageCount, book.PublishedYear, book.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM books WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
