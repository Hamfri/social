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
	User      User      `json:"user"`
}

type Comments interface {
	Create(context.Context, Comment) error
	GetCommentsByPostID(context.Context, int64) ([]*Comment, error)
}

type CommentRepository struct {
	DB *sql.DB
}

func (r CommentRepository) Create(ctx context.Context, comment Comment) error {
	return nil
}

func (r CommentRepository) GetCommentsByPostID(ctx context.Context, postId int64) ([]*Comment, error) {
	query := `
		SELECT c.user_id, u.id, u.username, c.id, c.comment, c.created_at 
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

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
