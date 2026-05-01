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

type CreateSupplierRequest struct {
	Name    string  `json:"name" binding:"required"`
	Code    *string `json:"code"`
	Contact *string `json:"contact"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`
	Status  *int    `json:"status"`
}

type UpdateSupplierRequest struct {
	Name    *string `json:"name"`
	Code    *string `json:"code"`
	Contact *string `json:"contact"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`
	Status  *int    `json:"status"`
}

type SupplierListResponse struct {
	Suppliers []model.Supplier `json:"items"`
	Total     int              `json:"total"`
	Page      int              `json:"page"`
	Size      int              `json:"size"`
}

type SupplierService interface {
	Create(ctx context.Context, input *service.CreateSupplierInput) (*model.Supplier, error)
	GetByID(ctx context.Context, id int64) (*model.Supplier, error)
	List(ctx context.Context, page, pageSize int, keyword string) (*service.ListSuppliersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateSupplierInput) (*model.Supplier, error)
	Delete(ctx context.Context, id int64) error
}

type SupplierHandler struct {
	supplierService SupplierService
}

func NewSupplierHandler(supplierService SupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

func (h *SupplierHandler) Create(c *gin.Context) {
	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateSupplierInput{
		Name: req.Name,
	}

	if req.Code != nil {
		input.Code = *req.Code
	}
	if req.Contact != nil {
		input.Contact = *req.Contact
	}
	if req.Phone != nil {
		input.Phone = *req.Phone
	}
	if req.Email != nil {
		input.Email = *req.Email
	}
	if req.Address != nil {
		input.Address = *req.Address
	}
	if req.Status != nil {
		input.Status = *req.Status
	}

	supplier, err := h.supplierService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleSupplierError(c, err)
		return
	}

	response.Success(c, supplier)
}

func (h *SupplierHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid supplier id")
		return
	}

	supplier, err := h.supplierService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleSupplierError(c, err)
		return
	}

	response.Success(c, supplier)
}

func (h *SupplierHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	result, err := h.supplierService.List(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		handleSupplierError(c, err)
		return
	}

	response.Success(c, SupplierListResponse{
		Suppliers: result.Suppliers,
		Total:     result.Total,
		Page:      page,
		Size:      pageSize,
	})
}

func (h *SupplierHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid supplier id")
		return
	}

	var req UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateSupplierInput{
		Name:    req.Name,
		Code:    req.Code,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
		Status:  req.Status,
	}

	supplier, err := h.supplierService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleSupplierError(c, err)
		return
	}

	response.Success(c, supplier)
}

func (h *SupplierHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid supplier id")
		return
	}

	err = h.supplierService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleSupplierError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleSupplierError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterSupplierRoutes(r *gin.RouterGroup, h *SupplierHandler) {
	suppliers := r.Group("/suppliers")
	{
		suppliers.GET("", h.List)
		suppliers.POST("", h.Create)
		suppliers.GET("/:id", h.GetByID)
		suppliers.PUT("/:id", h.Update)
		suppliers.DELETE("/:id", h.Delete)
	}
}
