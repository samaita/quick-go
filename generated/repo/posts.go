package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samaita/quick-go/generated/model"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) List(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, body, author_id, published, created_at, updated_at FROM posts LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("posts list: %w", err)
	}
	defer rows.Close()

	var items []*model.Post
	for rows.Next() {
		var m model.Post
		if err := rows.Scan(&m.Id, &m.Title, &m.Body, &m.AuthorId, &m.Published, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("posts scan: %w", err)
		}
		items = append(items, &m)
	}
	return items, rows.Err()
}

func (r *PostRepo) GetByID(ctx context.Context, id int64) (*model.Post, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, body, author_id, published, created_at, updated_at FROM posts WHERE id = ?`,
		id,
	)
	var m model.Post
	if err := row.Scan(&m.Id, &m.Title, &m.Body, &m.AuthorId, &m.Published, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("posts get: %w", err)
	}
	return &m, nil
}

func (r *PostRepo) Create(ctx context.Context, m *model.Post) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO posts (title, body, author_id, published, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		m.Title, m.Body, m.AuthorId, m.Published, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("posts create: %w", err)
	}
	return res.LastInsertId()
}

func (r *PostRepo) Update(ctx context.Context, id int64, m *model.Post) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE posts SET title = ?, body = ?, author_id = ?, published = ?, created_at = ?, updated_at = ? WHERE id = ?`,
		m.Title, m.Body, m.AuthorId, m.Published, m.CreatedAt, m.UpdatedAt, id,
	)
	if err != nil {
		return fmt.Errorf("posts update: %w", err)
	}
	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM posts WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("posts delete: %w", err)
	}
	return nil
}

// Suppress unused import lint warning when time.Time columns exist.
var _ = (*time.Time)(nil)
