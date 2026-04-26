package handler

import (
	"context"

	"warehouse/internal/middleware"
	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateThemeRequest struct {
	Theme string `json:"theme" binding:"required"`
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *model.User, error)
	GetProfile(ctx context.Context, userID int64) (*model.User, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	UpdateTheme(ctx context.Context, userID int64, theme string) error
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, apperrors.CodeUnauthorized, "user not authenticated")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, user)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, apperrors.CodeUnauthorized, "user not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *AuthHandler) UpdateTheme(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, apperrors.CodeUnauthorized, "user not authenticated")
		return
	}

	var req UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	err := h.authService.UpdateTheme(c.Request.Context(), userID, req.Theme)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleAuthError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}
