package mysql

import (
	"grab-bootcamp-be-team13-2025/internal/domain/models"

	"gorm.io/gorm"
	"context"
	"time"
)

type MySQLUserRepository struct {
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *MySQLUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MySQLUserRepository) UpdateUser(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *MySQLUserRepository) FindByResetToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MySQLUserRepository) SaveResetToken(ctx context.Context, email, token string, expiry time.Time) error {
    return r.db.WithContext(ctx).Model(&models.User{}).
        Where("email = ?", email).
        Updates(map[string]interface{}{
            "reset_token":       token,
            "reset_token_expiry": expiry,
        }).Error
}
