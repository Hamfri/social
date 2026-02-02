package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Posts interface {
	Create(context.Context, *Post) error
}

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type PostRepository struct {
	db *sql.DB
}

func (r *PostRepository) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts
		(content, title, user_id, tags)
		values($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	args := []any{post.Content, post.Title, post.UserID, pq.Array(post.Tags)}
	return r.db.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}
