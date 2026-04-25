package handler

import (
	"context"
	"strconv"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

type RoleListResponse struct {
	Roles []model.Role `json:"roles"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

type RoleService interface {
	List(ctx context.Context, page, pageSize int) ([]model.Role, int, error)
	GetByID(ctx context.Context, id int64) (*model.Role, error)
	Create(ctx context.Context, role *model.Role) (*model.Role, error)
	Update(ctx context.Context, id int64, role *model.Role) (*model.Role, error)
	Delete(ctx context.Context, id int64) error
}

type RoleHandler struct {
	roleService RoleService
}

func NewRoleHandler(roleService RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	roles, total, err := h.roleService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	response.Success(c, RoleListResponse{
		Roles: roles,
		Total: total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid role id")
		return
	}

	role, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	response.Success(c, role)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}

	createdRole, err := h.roleService.Create(c.Request.Context(), role)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	response.Success(c, createdRole)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid role id")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if req.Status != nil {
		role.Status = *req.Status
	}

	updatedRole, err := h.roleService.Update(c.Request.Context(), id, role)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	response.Success(c, updatedRole)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid role id")
		return
	}

	err = h.roleService.Delete(c.Request.Context(), id)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleRoleError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterRoleRoutes(r *gin.RouterGroup, h *RoleHandler) {
	roles := r.Group("/roles")
	{
		roles.GET("", h.List)
		roles.POST("", h.Create)
		roles.GET("/:id", h.GetByID)
		roles.PUT("/:id", h.Update)
		roles.DELETE("/:id", h.Delete)
	}
}
