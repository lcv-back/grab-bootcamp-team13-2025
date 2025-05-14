package models

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
