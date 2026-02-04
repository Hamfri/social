package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repository struct {
	Posts
	Users
	Comments
}

func New(db *sql.DB) Repository {
	return Repository{
		Posts:    &PostRepository{db},
		Users:    &UserRepository{db},
		Comments: &CommentRepository{db},
	}
}
