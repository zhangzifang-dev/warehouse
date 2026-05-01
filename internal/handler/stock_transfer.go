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

type CreateStockTransferRequest struct {
	OrderNo         string  `json:"order_no"`
	SourceWarehouseID int64   `json:"from_warehouse_id" binding:"required"`
	TargetWarehouseID   int64   `json:"to_warehouse_id" binding:"required"`
	TotalQty   float64 `json:"total_quantity"`
	Remark          string  `json:"remark"`
}

type UpdateStockTransferRequest struct {
	SourceWarehouseID *int64   `json:"from_warehouse_id,omitempty"`
	TargetWarehouseID   *int64   `json:"to_warehouse_id,omitempty"`
	TotalQty   *float64 `json:"total_quantity,omitempty"`
	Status          *int     `json:"status,omitempty"`
	Remark          *string  `json:"remark,omitempty"`
}

type StockTransferListResponse struct {
	Items []model.StockTransfer `json:"items"`
	Total int                   `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

type StockTransferService interface {
	Create(ctx context.Context, input *service.CreateStockTransferInput) (*model.StockTransfer, error)
	GetByID(ctx context.Context, id int64) (*model.StockTransfer, error)
	List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) (*service.ListStockTransfersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateStockTransferInput) (*model.StockTransfer, error)
	Delete(ctx context.Context, id int64) error
	Confirm(ctx context.Context, id int64) (*model.StockTransfer, error)
}

type StockTransferHandler struct {
	stockTransferService StockTransferService
}

func NewStockTransferHandler(stockTransferService StockTransferService) *StockTransferHandler {
	return &StockTransferHandler{
		stockTransferService: stockTransferService,
	}
}

func (h *StockTransferHandler) Create(c *gin.Context) {
	var req CreateStockTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateStockTransferInput{
		OrderNo:         req.OrderNo,
		SourceWarehouseID: req.SourceWarehouseID,
		TargetWarehouseID:   req.TargetWarehouseID,
		TotalQty:   req.TotalQty,
		Remark:          req.Remark,
	}

	transfer, err := h.stockTransferService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleStockTransferError(c, err)
		return
	}

	response.Success(c, transfer)
}

func (h *StockTransferHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid stock transfer id")
		return
	}

	transfer, err := h.stockTransferService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleStockTransferError(c, err)
		return
	}

	response.Success(c, transfer)
}

func (h *StockTransferHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	fromWarehouseID, _ := strconv.ParseInt(c.Query("from_warehouse_id"), 10, 64)
	toWarehouseID, _ := strconv.ParseInt(c.Query("to_warehouse_id"), 10, 64)
	
	status := -1
	if c.Query("status") != "" {
		status, _ = strconv.Atoi(c.Query("status"))
	}

	result, err := h.stockTransferService.List(c.Request.Context(), page, pageSize, int(fromWarehouseID), int(toWarehouseID), status)
	if err != nil {
		handleStockTransferError(c, err)
		return
	}

	response.Success(c, StockTransferListResponse{
		Items: result.Transfers,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *StockTransferHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid stock transfer id")
		return
	}

	var req UpdateStockTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateStockTransferInput{
		SourceWarehouseID: req.SourceWarehouseID,
		TargetWarehouseID:   req.TargetWarehouseID,
		TotalQty:   req.TotalQty,
		Status:          req.Status,
		Remark:          req.Remark,
	}

	transfer, err := h.stockTransferService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleStockTransferError(c, err)
		return
	}

	response.Success(c, transfer)
}

func (h *StockTransferHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid stock transfer id")
		return
	}

	err = h.stockTransferService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleStockTransferError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *StockTransferHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid stock transfer id")
		return
	}

	transfer, err := h.stockTransferService.Confirm(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleStockTransferError(c, err)
		return
	}

	response.Success(c, transfer)
}

func handleStockTransferError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterStockTransferRoutes(r *gin.RouterGroup, h *StockTransferHandler) {
	transfers := r.Group("/stock-transfers")
	{
		transfers.GET("", h.List)
		transfers.POST("", h.Create)
		transfers.GET("/:id", h.GetByID)
		transfers.PUT("/:id", h.Update)
		transfers.DELETE("/:id", h.Delete)
		transfers.POST("/:id/confirm", h.Confirm)
	}
}
