package mocks

import (
	"context"
	"errors"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is a mock type for repository.UserRepository
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) FindByResetToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) SaveResetToken(ctx context.Context, email, token string, expiry time.Time) error {
	args := m.Called(ctx, email, token, expiry)
	return args.Error(0)
}

// AuthUsecase struct for testing
type AuthUsecase struct {
	userRepo       *UserRepository
	jwtUtil        *JWTUtil
	redisClient    *RedisClient
	rabbitmqClient *RabbitMQClient
}

// Add the missing methods to implement the AuthUseCase interface
func (a *AuthUsecase) Signup(req *models.SignupRequest) error {
	existingUser, err := a.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	role := "user"
	if req.Role == "admin" {
		role = "admin"
	}

	user := &models.User{
		Email:    req.Email,
		Fullname: req.Fullname,
		Birthday: req.Birthday,
		Gender:   req.Gender,
		Password: string(hashedPassword),
		Role:     role,
	}

	return a.userRepo.CreateUser(user)
}

func (a *AuthUsecase) Login(req *models.LoginRequest) (string, error) {
	user, err := a.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := a.jwtUtil.GenerateToken(uint(user.ID), user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (a *AuthUsecase) GetUserInfo(email string) (*models.User, error) {
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (a *AuthUsecase) ForgotPassword(email string) error {
	return nil
}

func (a *AuthUsecase) ResetPassword(token, newPassword string) error {
	return nil
}

// NewAuthUsecase constructor for testing
func NewAuthUsecase(userRepo *UserRepository, jwtUtil *JWTUtil, redisClient *RedisClient, rabbitmqClient *RabbitMQClient) *AuthUsecase {
	return &AuthUsecase{
		userRepo:       userRepo,
		jwtUtil:        jwtUtil,
		redisClient:    redisClient,
		rabbitmqClient: rabbitmqClient,
	}
}
