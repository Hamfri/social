package repository

import (
	"context"
	"database/sql"
)

func MockNewRepository() Repository {
	return Repository{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	// users []User // append to mock redis
}

func (m *MockUserStore) GetByID(ctx context.Context, ID int64) (*User, error) {
	return &User{ID: 1}, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) Create(ctx context.Context, ts *sql.Tx, u *User) error {
	return nil
}

func (m *MockUserStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User) (*string, error) {
	return nil, nil
}

func (m *MockUserStore) GetUserByToken(ctx context.Context, tx *sql.Tx, scope, plainText string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) Activate(ctx context.Context, scope, plainText string) (*User, error) {
	return &User{}, nil
}
