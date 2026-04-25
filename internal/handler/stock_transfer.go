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

type CreateStockTransferRequest struct {
	OrderNo         string  `json:"order_no"`
	FromWarehouseID int64   `json:"from_warehouse_id" binding:"required"`
	ToWarehouseID   int64   `json:"to_warehouse_id" binding:"required"`
	TotalQuantity   float64 `json:"total_quantity"`
	Remark          string  `json:"remark"`
}

type UpdateStockTransferRequest struct {
	FromWarehouseID *int64   `json:"from_warehouse_id,omitempty"`
	ToWarehouseID   *int64   `json:"to_warehouse_id,omitempty"`
	TotalQuantity   *float64 `json:"total_quantity,omitempty"`
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
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		TotalQuantity:   req.TotalQuantity,
		Remark:          req.Remark,
	}

	transfer, err := h.stockTransferService.Create(c.Request.Context(), input)
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
	status, _ := strconv.Atoi(c.Query("status"))

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
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		TotalQuantity:   req.TotalQuantity,
		Status:          req.Status,
		Remark:          req.Remark,
	}

	transfer, err := h.stockTransferService.Update(c.Request.Context(), id, input)
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

	err = h.stockTransferService.Delete(c.Request.Context(), id)
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

	transfer, err := h.stockTransferService.Confirm(c.Request.Context(), id)
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
