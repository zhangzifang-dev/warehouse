package handler

import (
	"context"
	"strconv"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"
	"warehouse/internal/middleware"

	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  *int64 `json:"parent_id"`
	SortOrder *int   `json:"sort_order"`
	Status    *int   `json:"status"`
}

type UpdateCategoryRequest struct {
	Name      *string `json:"name"`
	ParentID  *int64  `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
	Status    *int    `json:"status"`
}

type CategoryListResponse struct {
	Categories []model.Category `json:"items"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
}

type CategoryService interface {
	Create(ctx context.Context, input *service.CreateCategoryInput) (*model.Category, error)
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	List(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateCategoryInput) (*model.Category, error)
	Delete(ctx context.Context, id int64) error
}

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateCategoryInput{
		Name: req.Name,
	}

	if req.ParentID != nil {
		input.ParentID = *req.ParentID
	}
	if req.SortOrder != nil {
		input.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		input.Status = *req.Status
	}

	category, err := h.categoryService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response.Success(c, category)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid category id")
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response.Success(c, category)
}

func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	name := c.Query("name")

	filter := &service.CategoryQueryFilter{
		Name:     name,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.categoryService.List(c.Request.Context(), filter)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response.Success(c, CategoryListResponse{
		Categories: result.Categories,
		Total:      result.Total,
		Page:       page,
		Size:       pageSize,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid category id")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateCategoryInput{
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	category, err := h.categoryService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response.Success(c, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid category id")
		return
	}

	err = h.categoryService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleCategoryError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleCategoryError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterCategoryRoutes(r *gin.RouterGroup, h *CategoryHandler) {
	categories := r.Group("/categories")
	{
		categories.GET("", h.List)
		categories.POST("", h.Create)
		categories.GET("/:id", h.GetByID)
		categories.PUT("/:id", h.Update)
		categories.DELETE("/:id", h.Delete)
	}
}
