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

type CreateInventoryRequest struct {
	WarehouseID int64    `json:"warehouse_id" binding:"required"`
	ProductID   int64    `json:"product_id" binding:"required"`
	LocationID  *int64   `json:"location_id"`
	Quantity    *float64 `json:"quantity"`
	BatchNo     *string  `json:"batch_no"`
}

type UpdateInventoryRequest struct {
	WarehouseID *int64   `json:"warehouse_id"`
	ProductID   *int64   `json:"product_id"`
	LocationID  *int64   `json:"location_id"`
	Quantity    *float64 `json:"quantity"`
	BatchNo     *string  `json:"batch_no"`
}

type AdjustQuantityRequest struct {
	InventoryID int64   `json:"inventory_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
}

type CheckStockRequest struct {
	WarehouseID int64   `json:"warehouse_id" binding:"required"`
	ProductID   int64   `json:"product_id" binding:"required"`
	BatchNo     string  `json:"batch_no"`
	Quantity    float64 `json:"quantity" binding:"required"`
}

type InventoryListResponse struct {
	Inventories []model.Inventory `json:"items"`
	Total       int               `json:"total"`
	Page        int               `json:"page"`
	Size        int               `json:"size"`
}

type CheckStockResponse struct {
	Available    bool    `json:"available"`
	CurrentStock float64 `json:"current_stock"`
	Requested    float64 `json:"requested"`
}

type InventoryService interface {
	Create(ctx context.Context, input *service.CreateInventoryInput) (*model.Inventory, error)
	GetByID(ctx context.Context, id int64) (*model.Inventory, error)
	List(ctx context.Context, filter *model.InventoryQueryFilter) (*service.ListInventoriesResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateInventoryInput) (*model.Inventory, error)
	Delete(ctx context.Context, id int64) error
	AdjustQuantity(ctx context.Context, input *service.AdjustQuantityInput) (*model.Inventory, error)
	CheckStock(ctx context.Context, input *service.CheckStockInput) (*service.CheckStockResult, error)
}

type InventoryHandler struct {
	inventoryService InventoryService
}

func NewInventoryHandler(inventoryService InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) Create(c *gin.Context) {
	var req CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateInventoryInput{
		WarehouseID: req.WarehouseID,
		ProductID:   req.ProductID,
	}

	if req.LocationID != nil {
		input.LocationID = *req.LocationID
	}
	if req.Quantity != nil {
		input.Quantity = *req.Quantity
	}
	if req.BatchNo != nil {
		input.BatchNo = *req.BatchNo
	}

	inventory, err := h.inventoryService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, inventory)
}

func (h *InventoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inventory id")
		return
	}

	inventory, err := h.inventoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, inventory)
}

func (h *InventoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	
	filter := &model.InventoryQueryFilter{
		Page:     page,
		PageSize: pageSize,
	}
	if warehouseID, err := strconv.ParseInt(c.Query("warehouse_id"), 10, 64); err == nil {
		filter.WarehouseID = warehouseID
	}

	if productID, err := strconv.ParseInt(c.Query("product_id"), 10, 64); err == nil {
		filter.ProductID = productID
	}
	
	if productName := c.Query("product_name"); productName != "" {
		filter.ProductName = productName
	}
	
	if batchNo := c.Query("batch_no"); batchNo != "" {
		filter.BatchNo = batchNo
	}
	
	if quantityMinStr := c.Query("quantity_min"); quantityMinStr != "" {
		if quantityMin, err := strconv.ParseFloat(quantityMinStr, 64); err == nil {
			filter.QuantityMin = &quantityMin
		}
	}
	
	if quantityMaxStr := c.Query("quantity_max"); quantityMaxStr != "" {
		if quantityMax, err := strconv.ParseFloat(quantityMaxStr, 64); err == nil {
			filter.QuantityMax = &quantityMax
		}
	}

	result, err := h.inventoryService.List(c.Request.Context(), filter)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, InventoryListResponse{
		Inventories: result.Inventories,
		Total:       result.Total,
		Page:        page,
		Size:        pageSize,
	})
}

func (h *InventoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inventory id")
		return
	}

	var req UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateInventoryInput{
		WarehouseID: req.WarehouseID,
		ProductID:   req.ProductID,
		LocationID:  req.LocationID,
		Quantity:    req.Quantity,
		BatchNo:     req.BatchNo,
	}

	inventory, err := h.inventoryService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, inventory)
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inventory id")
		return
	}

	err = h.inventoryService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *InventoryHandler) AdjustQuantity(c *gin.Context) {
	var req AdjustQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.AdjustQuantityInput{
		InventoryID: req.InventoryID,
		Quantity:    req.Quantity,
	}

	ctx := service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c))
	ctx = service.SetUserIDToContext(ctx, middleware.GetUserID(c))
	inventory, err := h.inventoryService.AdjustQuantity(ctx, input)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, inventory)
}

func (h *InventoryHandler) CheckStock(c *gin.Context) {
	var req CheckStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CheckStockInput{
		WarehouseID: req.WarehouseID,
		ProductID:   req.ProductID,
		BatchNo:     req.BatchNo,
		Quantity:    req.Quantity,
	}

	result, err := h.inventoryService.CheckStock(c.Request.Context(), input)
	if err != nil {
		handleInventoryError(c, err)
		return
	}

	response.Success(c, CheckStockResponse{
		Available:    result.Available,
		CurrentStock: result.CurrentStock,
		Requested:    result.Requested,
	})
}

func handleInventoryError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterInventoryRoutes(r *gin.RouterGroup, h *InventoryHandler) {
	inventories := r.Group("/inventory")
	{
		inventories.GET("", h.List)
		inventories.POST("", h.Create)
		inventories.GET("/:id", h.GetByID)
		inventories.PUT("/:id", h.Update)
		inventories.DELETE("/:id", h.Delete)
		inventories.POST("/adjust", h.AdjustQuantity)
		inventories.POST("/check", h.CheckStock)
	}
}
