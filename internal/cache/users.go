package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"social/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisCacheExp time.Duration = 2 * time.Hour

type Users interface {
	Get(context.Context, int64) (*repository.User, error)
	Set(context.Context, *repository.User) error
}

type UserStore struct {
	rdb *redis.Client
}

func (s *UserStore) Get(ctx context.Context, userId int64) (*repository.User, error) {
	if userId == 0 {
		return nil, errors.New("redis: invalid user id")
	}

	cacheKey := fmt.Sprintf("user-%d", userId)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var user repository.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *repository.User) error {
	if user.ID == 0 {
		return errors.New("redis: invalid user id")
	}
	cacheKey := fmt.Sprintf("user-%d", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	s.rdb.Set(ctx, cacheKey, json, redisCacheExp)
	return nil
}
