package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Repository struct {
	Posts
	Users
}

func New(db *sql.DB) Repository {
	return Repository{
		Posts: &PostRepository{db},
		Users: &UserRepository{db},
	}
}
