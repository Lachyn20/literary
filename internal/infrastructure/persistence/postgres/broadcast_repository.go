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

type BroadcastRepository struct {
	pool *pgxpool.Pool
}

func NewBroadcastRepository(pool *pgxpool.Pool) *BroadcastRepository {
	return &BroadcastRepository{pool: pool}
}

func (r *BroadcastRepository) Create(ctx context.Context, broadcast *entity.Broadcast) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO broadcasts (id, title, broadcast_type, channel_name, broadcast_date, file_path, file_type, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, broadcast.ID, broadcast.Title, broadcast.BroadcastType, broadcast.ChannelName, broadcast.BroadcastDate, broadcast.FilePath, broadcast.FileType, broadcast.CreatedAt)
	return err
}

func (r *BroadcastRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Broadcast, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,broadcast_type,channel_name,broadcast_date,file_path,file_type,created_at FROM broadcasts WHERE id=$1`, id)
	var broadcast entity.Broadcast
	if err := row.Scan(&broadcast.ID, &broadcast.Title, &broadcast.BroadcastType, &broadcast.ChannelName, &broadcast.BroadcastDate, &broadcast.FilePath, &broadcast.FileType, &broadcast.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &broadcast, nil
}

func (r *BroadcastRepository) List(ctx context.Context, limit, offset int) ([]*entity.Broadcast, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,title,broadcast_type,channel_name,broadcast_date,file_path,file_type,created_at FROM broadcasts ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var broadcasts []*entity.Broadcast
	for rows.Next() {
		var broadcast entity.Broadcast
		if err := rows.Scan(&broadcast.ID, &broadcast.Title, &broadcast.BroadcastType, &broadcast.ChannelName, &broadcast.BroadcastDate, &broadcast.FilePath, &broadcast.FileType, &broadcast.CreatedAt); err != nil {
			return nil, err
		}
		broadcasts = append(broadcasts, &broadcast)
	}
	return broadcasts, rows.Err()
}

func (r *BroadcastRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM broadcasts`).Scan(&count)
	return count, err
}

func (r *BroadcastRepository) Update(ctx context.Context, broadcast *entity.Broadcast) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE broadcasts SET title=$1, broadcast_type=$2, channel_name=$3, broadcast_date=$4, file_path=$5, file_type=$6 WHERE id=$7`, broadcast.Title, broadcast.BroadcastType, broadcast.ChannelName, broadcast.BroadcastDate, broadcast.FilePath, broadcast.FileType, broadcast.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *BroadcastRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM broadcasts WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
