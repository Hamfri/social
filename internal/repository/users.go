package repository

import (
	"context"
	"database/sql"
	"time"
)

type Users interface {
	Create(context.Context, *User) error
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users
		(username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	args := []any{user.Username, user.Email, user.Password}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
