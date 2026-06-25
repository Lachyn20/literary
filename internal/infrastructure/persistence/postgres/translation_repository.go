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

type TranslatedByAuthorRepository struct {
	pool *pgxpool.Pool
}

func NewTranslatedByAuthorRepository(pool *pgxpool.Pool) *TranslatedByAuthorRepository {
	return &TranslatedByAuthorRepository{pool: pool}
}

func (r *TranslatedByAuthorRepository) Create(ctx context.Context, translation *entity.TranslatedByAuthor) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO translated_by_author (id, original_author_name, original_language, work_title, notes) VALUES ($1,$2,$3,$4,$5)`, translation.ID, translation.OriginalAuthorName, translation.OriginalLanguage, translation.WorkTitle, translation.Notes)
	return err
}

func (r *TranslatedByAuthorRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TranslatedByAuthor, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,original_author_name,original_language,work_title,notes FROM translated_by_author WHERE id=$1`, id)
	var translation entity.TranslatedByAuthor
	if err := row.Scan(&translation.ID, &translation.OriginalAuthorName, &translation.OriginalLanguage, &translation.WorkTitle, &translation.Notes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &translation, nil
}

func (r *TranslatedByAuthorRepository) List(ctx context.Context) ([]*entity.TranslatedByAuthor, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,original_author_name,original_language,work_title,notes FROM translated_by_author`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var translations []*entity.TranslatedByAuthor
	for rows.Next() {
		var translation entity.TranslatedByAuthor
		if err := rows.Scan(&translation.ID, &translation.OriginalAuthorName, &translation.OriginalLanguage, &translation.WorkTitle, &translation.Notes); err != nil {
			return nil, err
		}
		translations = append(translations, &translation)
	}
	return translations, rows.Err()
}

func (r *TranslatedByAuthorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM translated_by_author WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

type TranslatedIntoLanguageRepository struct {
	pool *pgxpool.Pool
}

func NewTranslatedIntoLanguageRepository(pool *pgxpool.Pool) *TranslatedIntoLanguageRepository {
	return &TranslatedIntoLanguageRepository{pool: pool}
}

func (r *TranslatedIntoLanguageRepository) Create(ctx context.Context, translation *entity.TranslatedIntoLanguage) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO translated_into_languages (id, language_name, translator_name, work_title, notes) VALUES ($1,$2,$3,$4,$5)`, translation.ID, translation.LanguageName, translation.TranslatorName, translation.WorkTitle, translation.Notes)
	return err
}

func (r *TranslatedIntoLanguageRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TranslatedIntoLanguage, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,language_name,translator_name,work_title,notes FROM translated_into_languages WHERE id=$1`, id)
	var translation entity.TranslatedIntoLanguage
	if err := row.Scan(&translation.ID, &translation.LanguageName, &translation.TranslatorName, &translation.WorkTitle, &translation.Notes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &translation, nil
}

func (r *TranslatedIntoLanguageRepository) List(ctx context.Context) ([]*entity.TranslatedIntoLanguage, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,language_name,translator_name,work_title,notes FROM translated_into_languages`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var translations []*entity.TranslatedIntoLanguage
	for rows.Next() {
		var translation entity.TranslatedIntoLanguage
		if err := rows.Scan(&translation.ID, &translation.LanguageName, &translation.TranslatorName, &translation.WorkTitle, &translation.Notes); err != nil {
			return nil, err
		}
		translations = append(translations, &translation)
	}
	return translations, rows.Err()
}

func (r *TranslatedIntoLanguageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM translated_into_languages WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
