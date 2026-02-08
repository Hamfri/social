package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var (
	ErrAlreadyFollowing = errors.New("you are already following that user")
)

type UserFollow struct {
	FollowedID int64     `json:"followed_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserFollows interface {
	Follow(context.Context, *UserFollow) error
	Unfollow(context.Context, UserFollow) error
}

type UserFollowRepository struct {
	DB *sql.DB
}

func (r *UserFollowRepository) Follow(ctx context.Context, userFollow *UserFollow) error {
	query := `
		INSERT INTO user_follows
		(followed_id, follower_id)
		VALUES($1, $2)
		RETURNING created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	args := []any{userFollow.FollowedID, userFollow.FollowerID}

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&userFollow.CreatedAt,
	)

	var pqErr *pq.Error

	if err != nil {
		switch {
		// 23505 duplicate key value violates unique constraint
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return ErrAlreadyFollowing
		default:
			return err
		}
	}

	return nil
}

func (r *UserFollowRepository) Unfollow(ctx context.Context, userFollow UserFollow) error {
	query := `
		DELETE FROM user_follows
		WHERE followed_id = $1 and follower_id = $2 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	args := []any{userFollow.FollowedID, userFollow.FollowerID}

	results, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
