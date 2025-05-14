package interfaces

import "grab-bootcamp-be-team13-2025/internal/domain/models"

type AuthUseCase interface {
	Signup(req *models.SignupRequest) (*models.SignupResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetUserInfo(email string) (*models.User, error)
	ForgotPassword(email string) (*models.ForgotPasswordResponse, error)
	ResetPassword(token, newPassword string) (*models.ResetPasswordResponse, error)
	UpdateUserInfo(email string, req *models.UpdateUserRequest) (*models.User, error)
	VerifyPassword(hashedPassword string, plainPassword *string) bool
}
