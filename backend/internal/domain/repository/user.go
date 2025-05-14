package repository

import (
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"context"
	"time"
)

var ErrUserNotFound = "user not found"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	SaveResetToken(ctx context.Context, email, token string, expiry time.Time) error
	FindByResetToken(ctx context.Context, token string) (*models.User, error)
}
