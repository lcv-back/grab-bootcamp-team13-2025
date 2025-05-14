package mocks

import (
	"github.com/stretchr/testify/mock"
)

type JWTUtil struct {
	mock.Mock
}

func (m *JWTUtil) GenerateToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *JWTUtil) ValidateResetPasswordToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}