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

type CreateOutboundOrderRequest struct {
	OrderNo       string  `json:"order_no"`
	CustomerID    int64   `json:"customer_id"`
	WarehouseID   int64   `json:"warehouse_id" binding:"required"`
	TotalQuantity float64 `json:"total_quantity"`
	Remark        string  `json:"remark"`
}

type UpdateOutboundOrderRequest struct {
	CustomerID    *int64   `json:"customer_id,omitempty"`
	WarehouseID   *int64   `json:"warehouse_id,omitempty"`
	TotalQuantity *float64 `json:"total_quantity,omitempty"`
	Status        *int     `json:"status,omitempty"`
	Remark        *string  `json:"remark,omitempty"`
}

type OutboundOrderListResponse struct {
	Items []model.OutboundOrder `json:"items"`
	Total int                   `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

type OutboundOrderService interface {
	Create(ctx context.Context, input *service.CreateOutboundOrderInput) (*model.OutboundOrder, error)
	GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListOutboundOrdersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateOutboundOrderInput) (*model.OutboundOrder, error)
	Delete(ctx context.Context, id int64) error
	Confirm(ctx context.Context, id int64) (*model.OutboundOrder, error)
}

type OutboundOrderHandler struct {
	outboundOrderService OutboundOrderService
}

func NewOutboundOrderHandler(outboundOrderService OutboundOrderService) *OutboundOrderHandler {
	return &OutboundOrderHandler{
		outboundOrderService: outboundOrderService,
	}
}

func (h *OutboundOrderHandler) Create(c *gin.Context) {
	var req CreateOutboundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateOutboundOrderInput{
		OrderNo:       req.OrderNo,
		CustomerID:    req.CustomerID,
		WarehouseID:   req.WarehouseID,
		TotalQuantity: req.TotalQuantity,
		Remark:        req.Remark,
	}

	order, err := h.outboundOrderService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func (h *OutboundOrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid outbound order id")
		return
	}

	order, err := h.outboundOrderService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func (h *OutboundOrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	warehouseID, _ := strconv.ParseInt(c.Query("warehouse_id"), 10, 64)
	status, _ := strconv.Atoi(c.Query("status"))

	result, err := h.outboundOrderService.List(c.Request.Context(), page, pageSize, int(warehouseID), status)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, OutboundOrderListResponse{
		Items: result.Orders,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *OutboundOrderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid outbound order id")
		return
	}

	var req UpdateOutboundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateOutboundOrderInput{
		CustomerID:    req.CustomerID,
		WarehouseID:   req.WarehouseID,
		TotalQuantity: req.TotalQuantity,
		Status:        req.Status,
		Remark:        req.Remark,
	}

	order, err := h.outboundOrderService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func (h *OutboundOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid outbound order id")
		return
	}

	err = h.outboundOrderService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *OutboundOrderHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid outbound order id")
		return
	}

	order, err := h.outboundOrderService.Confirm(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleOutboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func handleOutboundOrderError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterOutboundOrderRoutes(r *gin.RouterGroup, h *OutboundOrderHandler) {
	orders := r.Group("/outbound-orders")
	{
		orders.GET("", h.List)
		orders.POST("", h.Create)
		orders.GET("/:id", h.GetByID)
		orders.PUT("/:id", h.Update)
		orders.DELETE("/:id", h.Delete)
		orders.POST("/:id/confirm", h.Confirm)
	}
}
