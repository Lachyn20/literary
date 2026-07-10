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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO users (id, name, email, password_hash, role, active, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`, user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.Active, user.CreatedAt)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,name,email,password_hash,role,active,created_at FROM users WHERE id=$1`, id)
	var user entity.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Active, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,name,email,password_hash,role,active,created_at FROM users WHERE email=$1`, email)
	var user entity.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Active, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,name,email,password_hash,role,active,created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Active, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE users SET name=$1, email=$2, role=$3, active=$4 WHERE id=$5`, user.Name, user.Email, user.Role, user.Active, user.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE users SET password_hash=$1 WHERE id=$2`, passwordHash, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
