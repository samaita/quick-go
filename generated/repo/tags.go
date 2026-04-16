package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samaita/quick-go/generated/model"
)

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) List(ctx context.Context, limit, offset int) ([]*model.Tag, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, slug, created_at FROM tags LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("tags list: %w", err)
	}
	defer rows.Close()

	var items []*model.Tag
	for rows.Next() {
		var m model.Tag
		if err := rows.Scan(&m.Id, &m.Name, &m.Slug, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("tags scan: %w", err)
		}
		items = append(items, &m)
	}
	return items, rows.Err()
}

func (r *TagRepo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at FROM tags WHERE id = ?`,
		id,
	)
	var m model.Tag
	if err := row.Scan(&m.Id, &m.Name, &m.Slug, &m.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("tags get: %w", err)
	}
	return &m, nil
}

func (r *TagRepo) Create(ctx context.Context, m *model.Tag) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO tags (name, slug, created_at) VALUES (?, ?, ?)`,
		m.Name, m.Slug, m.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("tags create: %w", err)
	}
	return res.LastInsertId()
}

func (r *TagRepo) Update(ctx context.Context, id int64, m *model.Tag) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tags SET name = ?, slug = ?, created_at = ? WHERE id = ?`,
		m.Name, m.Slug, m.CreatedAt, id,
	)
	if err != nil {
		return fmt.Errorf("tags update: %w", err)
	}
	return nil
}

func (r *TagRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM tags WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("tags delete: %w", err)
	}
	return nil
}

// Suppress unused import lint warning when time.Time columns exist.
var _ = (*time.Time)(nil)
