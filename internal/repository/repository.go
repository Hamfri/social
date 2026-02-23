package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrEditConflict      = errors.New("edit conflict")
	QueryTimeoutDuration = 5 * time.Minute
)

type Repository struct {
	Posts       Posts
	Users       Users
	Comments    Comments
	UserFollows UserFollows
	UserTokens  UserTokens
	Roles       Roles
}

func New(db *sql.DB) Repository {
	// avoid this in large codebases
	// can easily lead to spaghetti code
	// use services instead
	tokenRepo := &UserTokenRepository{db}
	roleRepo := &RoleRepository{db}
	userRepo := &UserRepository{db, tokenRepo, roleRepo}

	return Repository{
		Posts:       &PostRepository{db},
		Users:       userRepo,
		Comments:    &CommentRepository{db},
		UserFollows: &UserFollowRepository{db},
		UserTokens:  tokenRepo,
		Roles:       roleRepo,
	}
}

func WithTxAndResult[T any](ctx context.Context, db *sql.DB, fn func(*sql.Tx) (T, error)) (T, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		var val T
		return val, err
	}

	defer tx.Rollback()

	result, err := fn(tx)
	if err != nil {
		return result, err // rollback happens here via defer
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

func WithTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = fn(tx); err != nil {
		return err // rollback happens here via defer
	}

	return tx.Commit()
}
