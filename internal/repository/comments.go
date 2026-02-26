package repository

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	Comment   string    `json:"comment"`
	UserId    int64     `json:"user_id"`
	PostId    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	User      User      `json:"-"`
}

type Comments interface {
	Create(context.Context, *Comment) error
	GetCommentsByPostID(context.Context, int64) ([]*Comment, error)
}

type CommentRepository struct {
	DB *sql.DB
}

func (r CommentRepository) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments
		(comment, user_id, post_id)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	args := []any{comment.Comment, comment.UserId, comment.PostId}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
}

func (r CommentRepository) GetCommentsByPostID(ctx context.Context, postId int64) ([]*Comment, error) {
	query := `
		SELECT c.user_id, u.id, u.username, c.id, c.comment, c.created_at 
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// totalComments := 0
	comments := []*Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(
			&c.UserId,
			&c.User.ID,
			&c.User.Username,
			&c.ID,
			&c.Comment,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
