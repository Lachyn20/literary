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

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) VALUES ($1,$2,$3,$4,$5)`, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,user_id,token,expires_at,created_at FROM refresh_tokens WHERE token=$1`, token)
	var rt entity.RefreshToken
	if err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id=$1`, userID)
	return err
}
