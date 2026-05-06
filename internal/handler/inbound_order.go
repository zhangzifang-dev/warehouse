package handler

import (
	"context"
	"strconv"
	"time"

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
	ListWithFilter(ctx context.Context, filter *model.InboundOrderQueryFilter) (*service.ListInboundOrdersResult, error)
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

	filter := &model.InboundOrderQueryFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if orderNo := c.Query("order_no"); orderNo != "" {
		filter.OrderNo = orderNo
	}

	if supplierIDStr := c.Query("supplier_id"); supplierIDStr != "" {
		if supplierID, err := strconv.ParseInt(supplierIDStr, 10, 64); err == nil {
			filter.SupplierID = &supplierID
		}
	}

	if warehouseIDStr := c.Query("warehouse_id"); warehouseIDStr != "" {
		if warehouseID, err := strconv.ParseInt(warehouseIDStr, 10, 64); err == nil {
			filter.WarehouseID = &warehouseID
		}
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

	if createdAtStartStr := c.Query("created_at_start"); createdAtStartStr != "" {
		if createdAtStart, err := time.Parse(time.RFC3339, createdAtStartStr); err == nil {
			filter.CreatedAtStart = &createdAtStart
		}
	}

	if createdAtEndStr := c.Query("created_at_end"); createdAtEndStr != "" {
		if createdAtEnd, err := time.Parse(time.RFC3339, createdAtEndStr); err == nil {
			filter.CreatedAtEnd = &createdAtEnd
		}
	}

	result, err := h.inboundOrderService.ListWithFilter(c.Request.Context(), filter)
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
