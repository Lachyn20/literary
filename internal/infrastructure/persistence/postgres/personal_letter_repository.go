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

type PersonalLetterRepository struct {
	pool *pgxpool.Pool
}

func NewPersonalLetterRepository(pool *pgxpool.Pool) *PersonalLetterRepository {
	return &PersonalLetterRepository{pool: pool}
}

func (r *PersonalLetterRepository) Create(ctx context.Context, letter *entity.PersonalLetter) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO personal_letters (id, title, content, letter_date, scan_image_path, created_at) VALUES ($1,$2,$3,$4,$5,$6)`, letter.ID, letter.Title, letter.Content, letter.LetterDate, letter.ScanImagePath, letter.CreatedAt)
	return err
}

func (r *PersonalLetterRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PersonalLetter, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,content,letter_date,scan_image_path,created_at FROM personal_letters WHERE id=$1`, id)
	var letter entity.PersonalLetter
	if err := row.Scan(&letter.ID, &letter.Title, &letter.Content, &letter.LetterDate, &letter.ScanImagePath, &letter.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &letter, nil
}

func (r *PersonalLetterRepository) List(ctx context.Context) ([]*entity.PersonalLetter, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,content,letter_date,scan_image_path,created_at FROM personal_letters`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var letters []*entity.PersonalLetter
	for rows.Next() {
		var letter entity.PersonalLetter
		if err := rows.Scan(&letter.ID, &letter.Title, &letter.Content, &letter.LetterDate, &letter.ScanImagePath, &letter.CreatedAt); err != nil {
			return nil, err
		}
		letters = append(letters, &letter)
	}
	return letters, rows.Err()
}

func (r *PersonalLetterRepository) Update(ctx context.Context, letter *entity.PersonalLetter) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE personal_letters SET title=$1, content=$2, letter_date=$3, scan_image_path=$4 WHERE id=$5`, letter.Title, letter.Content, letter.LetterDate, letter.ScanImagePath, letter.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *PersonalLetterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM personal_letters WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
