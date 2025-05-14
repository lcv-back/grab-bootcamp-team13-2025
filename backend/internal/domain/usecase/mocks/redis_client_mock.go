package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type RedisClient struct {
	mock.Mock
}

func (m *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *RedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *RedisClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
    args := m.Called(ctx, key, expiration)
    return args.Error(0)
}

func (m *RedisClient) Incr(ctx context.Context, key string) error {
    args := m.Called(ctx, key)
    return args.Error(0)
}
