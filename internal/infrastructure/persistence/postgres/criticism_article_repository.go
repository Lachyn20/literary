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

type CriticismArticleRepository struct {
	pool *pgxpool.Pool
}

func NewCriticismArticleRepository(pool *pgxpool.Pool) *CriticismArticleRepository {
	return &CriticismArticleRepository{pool: pool}
}

func (r *CriticismArticleRepository) Create(ctx context.Context, article *entity.CriticismArticle) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO criticism_articles (id, title, content, author, publish_date) VALUES ($1,$2,$3,$4,$5)`, article.ID, article.Title, article.Content, article.Author, article.PublishDate)
	return err
}

func (r *CriticismArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CriticismArticle, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,content,author,publish_date FROM criticism_articles WHERE id=$1`, id)
	var article entity.CriticismArticle
	if err := row.Scan(&article.ID, &article.Title, &article.Content, &article.Author, &article.PublishDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &article, nil
}

func (r *CriticismArticleRepository) List(ctx context.Context) ([]*entity.CriticismArticle, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,content,author,publish_date FROM criticism_articles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*entity.CriticismArticle
	for rows.Next() {
		var article entity.CriticismArticle
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Author, &article.PublishDate); err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}
	return articles, rows.Err()
}

func (r *CriticismArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM criticism_articles WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
