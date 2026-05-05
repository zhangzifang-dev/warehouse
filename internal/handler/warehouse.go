package handler

import (
	"context"
	"strconv"

	"warehouse/internal/middleware"
	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateWarehouseRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Status  *int   `json:"status"`
}

type UpdateWarehouseRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Status  *int   `json:"status"`
}

type WarehouseListResponse struct {
	Items []model.Warehouse `json:"items"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

type WarehouseService interface {
	Create(ctx context.Context, input *service.CreateWarehouseInput) (*model.Warehouse, error)
	GetByID(ctx context.Context, id int64) (*model.Warehouse, error)
	List(ctx context.Context, filter *service.WarehouseFilter) (*service.ListWarehousesResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateWarehouseInput) (*model.Warehouse, error)
	Delete(ctx context.Context, id int64) error
}

type WarehouseHandler struct {
	warehouseService WarehouseService
}

func NewWarehouseHandler(warehouseService WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseService: warehouseService,
	}
}

func (h *WarehouseHandler) Create(c *gin.Context) {
	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateWarehouseInput{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		Contact: req.Contact,
		Phone:   req.Phone,
	}

	if req.Status != nil {
		input.Status = *req.Status
	}

	ctx := service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c))
	warehouse, err := h.warehouseService.Create(ctx, input)
	if err != nil {
		handleWarehouseError(c, err)
		return
	}

	response.Success(c, warehouse)
}

func (h *WarehouseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid warehouse id")
		return
	}

	warehouse, err := h.warehouseService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleWarehouseError(c, err)
		return
	}

	response.Success(c, warehouse)
}

func (h *WarehouseHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	name := c.Query("name")

	filter := &service.WarehouseFilter{
		Name:     name,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.warehouseService.List(c.Request.Context(), filter)
	if err != nil {
		handleWarehouseError(c, err)
		return
	}

	response.Success(c, WarehouseListResponse{
		Items: result.Warehouses,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *WarehouseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid warehouse id")
		return
	}

	var req UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateWarehouseInput{
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
		Phone:   req.Phone,
		Status:  req.Status,
	}

	ctx := service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c))
	ctx = service.SetUserIDToContext(ctx, middleware.GetUserID(c))
	warehouse, err := h.warehouseService.Update(ctx, id, input)
	if err != nil {
		handleWarehouseError(c, err)
		return
	}

	response.Success(c, warehouse)
}

func (h *WarehouseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid warehouse id")
		return
	}

	ctx := service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c))
	err = h.warehouseService.Delete(ctx, id)
	if err != nil {
		handleWarehouseError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleWarehouseError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterWarehouseRoutes(r *gin.RouterGroup, h *WarehouseHandler) {
	warehouses := r.Group("/warehouses")
	{
		warehouses.GET("", h.List)
		warehouses.POST("", h.Create)
		warehouses.GET("/:id", h.GetByID)
		warehouses.PUT("/:id", h.Update)
		warehouses.DELETE("/:id", h.Delete)
	}
}
