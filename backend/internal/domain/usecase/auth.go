package usecase

import (
	"context"
	"errors"
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/repository"
	"grab-bootcamp-be-team13-2025/internal/infrastructure/email"
	"log"
	"strconv"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo     repository.UserRepository
	jwtUtil      models.JWTUtil
	redisClient  models.RedisClient
	emailService email.EmailService
	bloomFilter  *bloom.BloomFilter
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtUtil models.JWTUtil, redisClient models.RedisClient, emailService email.EmailService, bloomFilter *bloom.BloomFilter) *AuthUsecase {
	return &AuthUsecase{
		userRepo:     userRepo,
		jwtUtil:      jwtUtil,
		redisClient:  redisClient,
		emailService: emailService,
		bloomFilter:  bloomFilter,
	}
}

func (a *AuthUsecase) Signup(req *models.SignupRequest) (*models.SignupResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.Fullname == "" {
		return nil, errors.New("missing required fields")
	}

	// Kiểm tra email trong Bloom Filter
	if a.bloomFilter.Test([]byte(req.Email)) {
		// Nếu Bloom Filter cho rằng "có thể tồn tại", truy vấn database để xác nhận
		existingUser, err := a.userRepo.GetByEmail(req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %v", err)
		}
		if existingUser != nil {
			return nil, fmt.Errorf("email already exists")
		}
	} else {
		log.Printf("Bloom Filter: email %s is definitely not in the set, skipping database query", req.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
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

	if err := a.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Thêm email vào Bloom Filter
	a.bloomFilter.Add([]byte(req.Email))
	log.Printf("Added email %s to Bloom Filter", req.Email)

	return &models.SignupResponse{
		Message: "user created successfully",
		Status:  "success",
	}, nil
}

func (a *AuthUsecase) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := a.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := a.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.LoginResponse{
		Message:     "login successful",
		Status:      "success",
		AccessToken: token,
	}, nil
}

func (a *AuthUsecase) ForgotPassword(email string) (*models.ForgotPasswordResponse, error) {
	// Rate limiting check
	ctx := context.Background()
	key := fmt.Sprintf("forgot_password_attempts:%s", email)
	attempts, err := a.redisClient.Get(ctx, key)
	if err == nil && attempts != "" {
		count, _ := strconv.Atoi(attempts)
		if count >= 10 {
			return nil, errors.New("too many attempts, please try again later")
		}
	}

	// Increment attempts
	err = a.redisClient.Incr(ctx, key)

	if err != nil {
		log.Printf("Redis INCR error: %v", err)
		return nil, errors.New("failed to track attempts")
	}
	// Set expiry for attempts counter (1 hour)
	err = a.redisClient.Expire(ctx, key, time.Hour)
	if err != nil {

		return nil, errors.New("failed to set attempts expiry")
	}

	// check if user exists
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// generate reset token
	resetToken, err := a.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("[ForgotPassword] failed to generate reset token: %v", err)
		return nil, errors.New("failed to generate reset token")
	}

	// save reset token to redis with time to live is 15 minutes
	err = a.redisClient.Set(ctx, "reset_password: "+resetToken, email, 15*time.Minute)
	if err != nil {
		log.Printf("[ForgotPassword] failed to save reset token to redis: %v", err)
		return nil, errors.New("failed to save reset token to redis")
	}

	// generate reset link using env or fallback
	resetURL := "https://isymptom.vercel.app/reset-password"
	resetLink := resetURL + "?token=" + resetToken

	// Gửi email trực tiếp
	err = a.emailService.SendResetPasswordEmail(ctx, user.Fullname, user.Email, resetLink)
	if err != nil {
		log.Printf("[ForgotPassword] failed to send email: %v", err)
		return nil, errors.New("failed to send reset password email")
	}

	// Log the forgot password request
	log.Printf("Forgot password request for email: %s", email)

	return &models.ForgotPasswordResponse{
		Message: "reset password email sent",
		Status:  "success",
	}, nil
}

func (a *AuthUsecase) ResetPassword(token, newPassword string) (*models.ResetPasswordResponse, error) {
	// check if reset token exists in redis
	ctx := context.Background()
	email, err := a.redisClient.Get(ctx, "reset_password: "+token)
	if err != nil {
		return nil, errors.New("invalid or expired reset token")
	}

	// delete token from redis
	err = a.redisClient.Delete(ctx, "reset_password: "+token)
	if err != nil {
		return nil, errors.New("failed to invalidate reset token")
	}

	// validate reset token
	validateEmail, err := a.jwtUtil.ValidateResetPasswordToken(token)
	if err != nil || validateEmail != email {
		return nil, errors.New("invalid or expired reset token")
	}

	if email != validateEmail {
		return nil, errors.New("invalid or expired reset token")
	}

	// find user by email
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash new password")
	}

	// update user password
	user.Password = string(hashedPassword)
	err = a.userRepo.UpdateUser(user)
	if err != nil {
		return nil, errors.New("failed to update password")
	}

	return &models.ResetPasswordResponse{
		Message: "password reset successful",
		Status:  "success",
	}, nil
}

func (a *AuthUsecase) UpdateUserInfo(email string, req *models.UpdateUserRequest) (*models.User, error) {
	// Lấy thông tin user hiện tại
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Cập nhật thông tin cơ bản nếu được cung cấp
	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}
	if req.Birthday != nil {
		user.Birthday = *req.Birthday
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}

	// Nếu có yêu cầu đổi mật khẩu
	if req.NewPassword != nil {
		// Hash mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash new password")
		}
		user.Password = string(hashedPassword)
	}

	// Lưu vào database
	if err := a.userRepo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update user information")
	}

	return user, nil
}

func (a *AuthUsecase) GetUserInfo(email string) (*models.User, error) {
	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (a *AuthUsecase) VerifyPassword(hashedPassword string, plainPassword *string) bool {
	if plainPassword == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(*plainPassword))
	return err == nil
}
