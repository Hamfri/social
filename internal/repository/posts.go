package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social/internal/pagination"
	"time"

	"github.com/lib/pq"
)

type Posts interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	DeletePost(context.Context, int64) error
	UpdatePost(context.Context, *Post) error
	GetUserFeed(context.Context, int64, *pagination.Pagination, pagination.Filter) ([]*PostWithMetadata, error)
}

type Post struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	Version   int        `json:"version"`
	UserID    int64      `json:"user_id"`
	Comments  []*Comment `json:"comments"` // should move to a dto or a domain
	User      User       `json:"user"`     // should move to a dto or a domain
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
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
	args := []any{
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (r *PostRepository) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, title, content, tags, user_id, version
		FROM posts
		WHERE id = $1
	`

	args := []any{id}

	var post Post

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.UserID,
		&post.Version,
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
		SET title = $1, content = $2, tags = $3, version = version + 1
		WHERE id = $4 and version = $5
		RETURNING updated_at, version
	`
	args := []any{
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
		post.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.UpdatedAt,
		&post.Version,
	)
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

// followedId equals the the authenticated user
func (r PostRepository) GetUserFeed(ctx context.Context, followedID int64, pagination *pagination.Pagination, filter pagination.Filter) ([]*PostWithMetadata, error) {
	query := `
		WITH related_users AS (
		    -- users followed by the authenticated user
			SELECT followed_id AS user_id
			FROM user_follows
			WHERE follower_id = $1

			UNION

			-- users following authenticated user
			SELECT follower_id AS user_id
			FROM user_follows
			WHERE followed_id = $1

			UNION

			-- we need to include the authenticated user_id
			-- as well so that we can return his/her posts too
			SELECT $1 AS user_id
		) 
		SELECT COUNT(*) OVER(), p.id as post_id, p.title, p.content, p.tags, u.id, u.username, c.comments_count
		FROM related_users ru
		JOIN users u ON u.id = ru.user_id
		JOIN posts p ON p.user_id = u.id

		-- Optimized incase comments table get's large
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comments_count
			FROM comments
			GROUP BY post_id
		) c ON c.post_id = p.id

		-- '||' concatenation operator
		WHERE (p.title || ' ' || p.content) ILIKE '%' || $4 || '%'

		-- if tags are provided filter by them if not don't filter
		-- type cast to text '::text[]'
		AND (p.tags @> $5 OR $5 = '{}'::text[])
		ORDER BY p.` + fmt.Sprintf(`%s %s`, filter.SortColumn(), filter.SortDirection()) + ` LIMIT $2 OFFSET $3`

	args := []any{
		followedID,
		pagination.PageSize,
		pagination.Offset(),
		filter.Search,
		pq.Array(filter.Tags),
	}
	fmt.Println(filter.Search)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	totalRecords := 0
	var postsWithMetadata []*PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata

		err := rows.Scan(
			&totalRecords,
			&p.ID,
			&p.Title,
			&p.Content,
			pq.Array(&p.Tags),
			&p.User.ID,
			&p.User.Username,
			&p.CommentsCount,
		)

		if err != nil {
			return nil, err
		}

		postsWithMetadata = append(postsWithMetadata, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pagination.TotalRecords = totalRecords

	return postsWithMetadata, nil
}
