package postgres

import (
    "context"

    "github.com/google/uuid"
    "github.com/hemra-siirow/literary/internal/domain/entity"
    "github.com/hemra-siirow/literary/internal/domain/repository"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type BiographyEventRepository struct {
    pool *pgxpool.Pool
}

func NewBiographyEventRepository(pool *pgxpool.Pool) *BiographyEventRepository {
    return &BiographyEventRepository{pool: pool}
}

func (r *BiographyEventRepository) Create(ctx context.Context, event *entity.BiographyEvent) error {
    _, err := r.pool.Exec(ctx, `INSERT INTO biography_events (id, biography_id, year, title_tk, title_ru, title_en, description_tk, description_ru, description_en, sort_order, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
        event.ID, event.BiographyID, event.Year, event.TitleTk, event.TitleRu, event.TitleEn, event.DescriptionTk, event.DescriptionRu, event.DescriptionEn, event.SortOrder, event.CreatedAt)
    return err
}

func (r *BiographyEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.BiographyEvent, error) {
    row := r.pool.QueryRow(ctx, `SELECT id, biography_id, year, title_tk, title_ru, title_en, description_tk, description_ru, description_en, sort_order, created_at FROM biography_events WHERE id=$1`, id)
    var e entity.BiographyEvent
    if err := row.Scan(&e.ID, &e.BiographyID, &e.Year, &e.TitleTk, &e.TitleRu, &e.TitleEn, &e.DescriptionTk, &e.DescriptionRu, &e.DescriptionEn, &e.SortOrder, &e.CreatedAt); err != nil {
        if err == pgx.ErrNoRows {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }
    return &e, nil
}

func (r *BiographyEventRepository) ListByBiographyID(ctx context.Context, biographyID uuid.UUID) ([]*entity.BiographyEvent, error) {
    rows, err := r.pool.Query(ctx, `SELECT id, biography_id, year, title_tk, title_ru, title_en, description_tk, description_ru, description_en, sort_order, created_at FROM biography_events WHERE biography_id=$1 ORDER BY year ASC, sort_order ASC`, biographyID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    events := make([]*entity.BiographyEvent, 0)
    for rows.Next() {
        var e entity.BiographyEvent
        if err := rows.Scan(&e.ID, &e.BiographyID, &e.Year, &e.TitleTk, &e.TitleRu, &e.TitleEn, &e.DescriptionTk, &e.DescriptionRu, &e.DescriptionEn, &e.SortOrder, &e.CreatedAt); err != nil {
            return nil, err
        }
        ev := e
        events = append(events, &ev)
    }
    return events, nil
}

func (r *BiographyEventRepository) Update(ctx context.Context, event *entity.BiographyEvent) error {
    cmd, err := r.pool.Exec(ctx, `UPDATE biography_events SET year=$1, title_tk=$2, title_ru=$3, title_en=$4, description_tk=$5, description_ru=$6, description_en=$7, sort_order=$8 WHERE id=$9`,
        event.Year, event.TitleTk, event.TitleRu, event.TitleEn, event.DescriptionTk, event.DescriptionRu, event.DescriptionEn, event.SortOrder, event.ID)
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return repository.ErrNotFound
    }
    return nil
}

func (r *BiographyEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
    _, err := r.pool.Exec(ctx, `DELETE FROM biography_events WHERE id=$1`, id)
    return err
}
