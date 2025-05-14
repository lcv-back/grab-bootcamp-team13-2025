package handlers

import (
	"grab-bootcamp-be-team13-2025/internal/domain/interfaces"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase interfaces.AuthUseCase
}

func NewAuthHandler(authUsecase interfaces.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "data invalid: " + err.Error(),
		})
		return
	}

	if _, err := h.authUsecase.Signup(&req); err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, models.Response{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SignupResponse{
		Message: "create user successfully",
		Status:  "success",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "data invalid: " + err.Error(),
		})
		return
	}

	loginResp, err := h.authUsecase.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginResp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	// Lấy email từ token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Error: "unauthorized",
		})
		return
	}

	// convert claims to map
	claimsMap := claims.(jwt.MapClaims)
	email := claimsMap["email"].(string)

	// get information user
	user, err := h.authUsecase.GetUserInfo(email)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "invalid email format",
		})
		return
	}

	if _, err := h.authUsecase.ForgotPassword(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "reset password link has been sent to your email",
		Status:  "success",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "invalid request format",
		})
		return
	}

	if _, err := h.authUsecase.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "password has been reset successfully",
		Status:  "success",
	})
}

func (h *AuthHandler) UpdateInfo(c *gin.Context) {
	// Lấy email từ token JWT
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Error: "unauthorized",
		})
		return
	}

	// Convert claims to map
	claimsMap := claims.(jwt.MapClaims)
	email := claimsMap["email"].(string)

	// Parse request body
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "data invalid: " + err.Error(),
		})
		return
	}

	// Kiểm tra nếu có yêu cầu đổi mật khẩu
	if req.NewPassword != nil {
		if req.CurrentPassword == nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Error: "current password is required when changing password",
			})
			return
		}

		// Kiểm tra mật khẩu hiện tại
		currentUser, err := h.authUsecase.GetUserInfo(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Error: "failed to get user information",
			})
			return
		}

		// Verify mật khẩu hiện tại
		if !h.authUsecase.VerifyPassword(currentUser.Password, req.CurrentPassword) {
			c.JSON(http.StatusBadRequest, models.Response{
				Error: "current password is incorrect",
			})
			return
		}
	}

	// Cập nhật thông tin user
	_, err := h.authUsecase.UpdateUserInfo(email, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Message: "user information updated successfully",
		Status:  "success",
	})
}
