package repository

import (
	"context"
	"database/sql"
	"errors"
)

type Roles interface {
	GetByName(ctx context.Context, name string) (*Role, error)
}

type Role struct {
	ID          int64
	Name        string
	Level       int64
	Description string
}

type RoleRepository struct {
	DB *sql.DB
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, level, description 
		FROM roles
		WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var role Role
	err := r.DB.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}
