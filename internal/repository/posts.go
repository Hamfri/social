package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Posts interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	DeletePost(context.Context, int64) error
	UpdatePost(context.Context, *Post) error
}

type Post struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	UserID    int64      `json:"user_id"`
	Comments  []*Comment `json:"comments"`
}

type PostRepository struct {
	DB *sql.DB
}

func (r *PostRepository) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts
		(content, title, user_id, tags)
		values($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	args := []any{post.Content, post.Title, post.UserID, pq.Array(post.Tags)}
	return r.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *PostRepository) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, title, content, tags, user_id
		FROM posts
		WHERE id = $1
	`

	args := []any{id}

	var post Post
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, tags = $3
		WHERE id = $4
		RETURNING updated_at
	`
	args := []any{
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
	}

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&post.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *PostRepository) DeletePost(ctx context.Context, postId int64) error {
	query := `
		DELETE FROM posts 
		WHERE id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
