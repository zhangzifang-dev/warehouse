package handler

import (
	"context"
	"strconv"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   *int   `json:"status"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   *int   `json:"status"`
}

type UserListResponse struct {
	Users []model.User `json:"users"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

type UserService interface {
	Create(ctx context.Context, input *service.CreateUserInput) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context, page, pageSize int) (*service.ListUsersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if req.Status != nil {
		input.Status = *req.Status
	}

	user, err := h.userService.Create(c.Request.Context(), input)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result, err := h.userService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, UserListResponse{
		Users: result.Users,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateUserInput{
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}

	user, err := h.userService.Update(c.Request.Context(), id, input)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid user id")
		return
	}

	err = h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid user id")
		return
	}

	roles, err := h.userService.GetUserRoles(c.Request.Context(), id)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response.Success(c, roles)
}

func handleUserError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterUserRoutes(r *gin.RouterGroup, h *UserHandler) {
	users := r.Group("/users")
	{
		users.GET("", h.List)
		users.POST("", h.Create)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
		users.GET("/:id/roles", h.GetRoles)
	}
}
