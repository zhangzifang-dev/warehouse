package handler

import (
	"context"
	"strconv"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type PermissionListResponse struct {
	Items []model.Permission `json:"items"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

type PermissionService interface {
	List(ctx context.Context, page, pageSize int) ([]model.Permission, int, error)
	GetByID(ctx context.Context, id int64) (*model.Permission, error)
}

type PermissionHandler struct {
	permService PermissionService
}

func NewPermissionHandler(permService PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permService: permService,
	}
}

func (h *PermissionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	permissions, total, err := h.permService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, apperrors.CodeInternalError, "failed to list permissions")
		return
	}

	response.Success(c, PermissionListResponse{
		Items: permissions,
		Total: total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *PermissionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid permission id")
		return
	}

	permission, err := h.permService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, apperrors.CodeNotFound, "permission not found")
		return
	}

	response.Success(c, permission)
}

func RegisterPermissionRoutes(r *gin.RouterGroup, h *PermissionHandler) {
	permissions := r.Group("/permissions")
	{
		permissions.GET("", h.List)
		permissions.GET("/:id", h.GetByID)
	}
}
