package cache

import (
	"context"
	"social/internal/repository"
)

func MockNewRedisStorage() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m MockUserStore) Get(ctx context.Context, userId int64) (*repository.User, error) {
	return &repository.User{}, nil
}

func (m MockUserStore) Set(ctx context.Context, user *repository.User) error {
	return nil
}
