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

type CreateInboundOrderRequest struct {
	OrderNo       string  `json:"order_no"`
	SupplierID    int64   `json:"supplier_id"`
	WarehouseID   int64   `json:"warehouse_id" binding:"required"`
	TotalQuantity float64 `json:"total_quantity"`
	Remark        string  `json:"remark"`
}

type UpdateInboundOrderRequest struct {
	SupplierID    *int64   `json:"supplier_id,omitempty"`
	WarehouseID   *int64   `json:"warehouse_id,omitempty"`
	TotalQuantity *float64 `json:"total_quantity,omitempty"`
	Status        *int     `json:"status,omitempty"`
	Remark        *string  `json:"remark,omitempty"`
}

type InboundOrderListResponse struct {
	Items []model.InboundOrder `json:"items"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

type InboundOrderService interface {
	Create(ctx context.Context, input *service.CreateInboundOrderInput) (*model.InboundOrder, error)
	GetByID(ctx context.Context, id int64) (*model.InboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListInboundOrdersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateInboundOrderInput) (*model.InboundOrder, error)
	Delete(ctx context.Context, id int64) error
	Confirm(ctx context.Context, id int64) (*model.InboundOrder, error)
}

type InboundOrderHandler struct {
	inboundOrderService InboundOrderService
}

func NewInboundOrderHandler(inboundOrderService InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{
		inboundOrderService: inboundOrderService,
	}
}

func (h *InboundOrderHandler) Create(c *gin.Context) {
	var req CreateInboundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateInboundOrderInput{
		OrderNo:       req.OrderNo,
		SupplierID:    req.SupplierID,
		WarehouseID:   req.WarehouseID,
		TotalQuantity: req.TotalQuantity,
		Remark:        req.Remark,
	}

	order, err := h.inboundOrderService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleInboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func (h *InboundOrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inbound order id")
		return
	}

	order, err := h.inboundOrderService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleInboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func (h *InboundOrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	warehouseID, _ := strconv.ParseInt(c.Query("warehouse_id"), 10, 64)
	status, _ := strconv.Atoi(c.Query("status"))

	result, err := h.inboundOrderService.List(c.Request.Context(), page, pageSize, int(warehouseID), status)
	if err != nil {
		handleInboundOrderError(c, err)
		return
	}

	response.Success(c, InboundOrderListResponse{
		Items: result.Orders,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *InboundOrderHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inbound order id")
		return
	}

	var req UpdateInboundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateInboundOrderInput{
		SupplierID:    req.SupplierID,
		WarehouseID:   req.WarehouseID,
		TotalQuantity: req.TotalQuantity,
		Status:        req.Status,
		Remark:        req.Remark,
	}

	order, err := h.inboundOrderService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleInboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func (h *InboundOrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inbound order id")
		return
	}

	err = h.inboundOrderService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleInboundOrderError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *InboundOrderHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid inbound order id")
		return
	}

	order, err := h.inboundOrderService.Confirm(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleInboundOrderError(c, err)
		return
	}

	response.Success(c, order)
}

func handleInboundOrderError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterInboundOrderRoutes(r *gin.RouterGroup, h *InboundOrderHandler) {
	orders := r.Group("/inbound-orders")
	{
		orders.GET("", h.List)
		orders.POST("", h.Create)
		orders.GET("/:id", h.GetByID)
		orders.PUT("/:id", h.Update)
		orders.DELETE("/:id", h.Delete)
		orders.POST("/:id/confirm", h.Confirm)
	}
}
