package repository

import (
	"database/sql"
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
