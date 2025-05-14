package models

import "time"

type PasswordResetToken struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id"`
	Token     string    `json:"token" gorm:"type:varchar(255);index:idx_token,length:255"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
