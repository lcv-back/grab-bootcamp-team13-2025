package models

type JWTUtil interface {
	GenerateToken(id uint, email string) (string, error)
	ValidateResetPasswordToken(token string) (string, error)
}
