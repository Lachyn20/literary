package postgres

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkRepository struct {
	pool *pgxpool.Pool
}

func NewWorkRepository(pool *pgxpool.Pool) *WorkRepository {
	return &WorkRepository{pool: pool}
}

func (r *WorkRepository) Create(ctx context.Context, work *entity.Work) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO works (id, title, category_id, file_path, description, audience_type, publish_year, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, work.ID, work.Title, work.CategoryID, work.FilePath, work.Description, work.AudienceType, work.PublishYear, work.CreatedAt, work.UpdatedAt)
	return err
}

func (r *WorkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	row := r.pool.QueryRow(ctx, `SELECT id,title,category_id,file_path,description,audience_type,publish_year,created_at,updated_at FROM works WHERE id=$1`, id)
	var work entity.Work
	if err := row.Scan(&work.ID, &work.Title, &work.CategoryID, &work.FilePath, &work.Description, &work.AudienceType, &work.PublishYear, &work.CreatedAt, &work.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &work, nil
}

func (r *WorkRepository) List(ctx context.Context, filter repository.WorkFilter, limit, offset int) ([]*entity.Work, error) {
	query := `SELECT id,title,category_id,file_path,description,audience_type,publish_year,created_at,updated_at FROM works`
	args := []interface{}{}
	clauses := []string{}
	if filter.CategoryID != nil {
		clauses = append(clauses, `category_id=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.CategoryID)
	}
	if filter.AudienceType != nil {
		clauses = append(clauses, `audience_type=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.AudienceType)
	}
	if filter.PublishYear != nil {
		clauses = append(clauses, `publish_year=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.PublishYear)
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	
	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []*entity.Work
	for rows.Next() {
		var work entity.Work
		if err := rows.Scan(&work.ID, &work.Title, &work.CategoryID, &work.FilePath, &work.Description, &work.AudienceType, &work.PublishYear, &work.CreatedAt, &work.UpdatedAt); err != nil {
			return nil, err
		}
		works = append(works, &work)
	}
	return works, rows.Err()
}

func (r *WorkRepository) Count(ctx context.Context, filter repository.WorkFilter) (int, error) {
	query := `SELECT COUNT(*) FROM works`
	args := []interface{}{}
	clauses := []string{}
	if filter.CategoryID != nil {
		clauses = append(clauses, `category_id=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.CategoryID)
	}
	if filter.AudienceType != nil {
		clauses = append(clauses, `audience_type=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.AudienceType)
	}
	if filter.PublishYear != nil {
		clauses = append(clauses, `publish_year=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.PublishYear)
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	
	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *WorkRepository) Update(ctx context.Context, work *entity.Work) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE works SET title=$1, category_id=$2, file_path=$3, description=$4, audience_type=$5, publish_year=$6, updated_at=$7 WHERE id=$8`, work.Title, work.CategoryID, work.FilePath, work.Description, work.AudienceType, work.PublishYear, work.UpdatedAt, work.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *WorkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM works WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *WorkRepository) Search(ctx context.Context, filter repository.WorkFilter) ([]*entity.Work, int, error) {
	// Build base where clauses similar to List
	where := []string{}
	args := []interface{}{}
	if filter.CategoryID != nil {
		where = append(where, `category_id=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.CategoryID)
	}
	if filter.AudienceType != nil {
		where = append(where, `audience_type=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.AudienceType)
	}
	if filter.PublishYear != nil {
		where = append(where, `publish_year=$`+strconv.Itoa(len(args)+1))
		args = append(args, *filter.PublishYear)
	}

	// full-text search
	if filter.Search != nil && *filter.Search != "" {
		where = append(where, `search_vector @@ plainto_tsquery('simple', $`+strconv.Itoa(len(args)+1)+`)`)
		args = append(args, *filter.Search)
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = " WHERE " + strings.Join(where, " AND ")
	}

	// count total
	countQuery := `SELECT COUNT(1) FROM works` + whereSQL
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	// select with rank ordering when searching
	query := `SELECT id,title,category_id,file_path,description,audience_type,publish_year,created_at,updated_at FROM works` + whereSQL
	if filter.Search != nil && *filter.Search != "" {
		// prefer ordering by rank
		query = `SELECT id,title,category_id,file_path,description,audience_type,publish_year,created_at,updated_at FROM works` + whereSQL + ` ORDER BY ts_rank(search_vector, plainto_tsquery('simple', $` + strconv.Itoa(len(args)) + `)) DESC` + ` LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		args = append(args, limit, offset)
	} else {
		query = query + ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		args = append(args, limit, offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var works []*entity.Work
	for rows.Next() {
		var work entity.Work
		if err := rows.Scan(&work.ID, &work.Title, &work.CategoryID, &work.FilePath, &work.Description, &work.AudienceType, &work.PublishYear, &work.CreatedAt, &work.UpdatedAt); err != nil {
			return nil, 0, err
		}
		works = append(works, &work)
	}
	return works, total, rows.Err()
}
