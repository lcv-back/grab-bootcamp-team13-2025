// internal/domain/models/user.go
package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Fullname  string    `json:"fullname" gorm:"not null"`
	Birthday  time.Time `json:"birthday" gorm:"not null" `
	Gender    string    `json:"gender" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:user"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SignupRequest struct {
	Email    string    `json:"email" binding:"required,email"`
	Fullname string    `json:"fullname" binding:"required,min=3,max=100"`
	Birthday time.Time `json:"birthday" binding:"required"`
	Gender   string    `json:"gender" binding:"required,oneof=male female other"`
	Password string    `json:"password" binding:"required,min=8,max=100"`
	Role     string    `json:"role" binding:"omitempty,oneof=user admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=100"`
}

type SignupResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type LoginResponse struct {
	Message     string `json:"message"`
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Response struct {
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UpdateUserRequest struct {
	Fullname        *string    `json:"fullname" binding:"omitempty,min=3,max=100"`
	Birthday        *time.Time `json:"birthday" binding:"omitempty"`
	Gender          *string    `json:"gender" binding:"omitempty,oneof=male female other"`
	CurrentPassword *string    `json:"current_password" binding:"omitempty,min=8,max=100"`
	NewPassword     *string    `json:"new_password" binding:"omitempty,min=8,max=100,required_with=CurrentPassword"`
}
