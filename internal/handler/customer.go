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

type CreateCustomerRequest struct {
	Name    string  `json:"name" binding:"required"`
	Code    *string `json:"code"`
	Contact *string `json:"contact"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`
	Status  *int    `json:"status"`
}

type UpdateCustomerRequest struct {
	Name    *string `json:"name"`
	Code    *string `json:"code"`
	Contact *string `json:"contact"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`
	Status  *int    `json:"status"`
}

type CustomerListResponse struct {
	Customers []model.Customer `json:"items"`
	Total     int              `json:"total"`
	Page      int              `json:"page"`
	Size      int              `json:"size"`
}

type CustomerService interface {
	Create(ctx context.Context, input *service.CreateCustomerInput) (*model.Customer, error)
	GetByID(ctx context.Context, id int64) (*model.Customer, error)
	List(ctx context.Context, filter *service.CustomerQueryFilter) (*service.ListCustomersResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateCustomerInput) (*model.Customer, error)
	Delete(ctx context.Context, id int64) error
}

type CustomerHandler struct {
	customerService CustomerService
}

func NewCustomerHandler(customerService CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateCustomerInput{
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

	customer, err := h.customerService.Create(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), input)
	if err != nil {
		handleCustomerError(c, err)
		return
	}

	response.Success(c, customer)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid customer id")
		return
	}

	customer, err := h.customerService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleCustomerError(c, err)
		return
	}

	response.Success(c, customer)
}

func (h *CustomerHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	filter := &service.CustomerQueryFilter{
		Code:     c.Query("code"),
		Name:     c.Query("name"),
		Phone:    c.Query("phone"),
		Page:     page,
		PageSize: pageSize,
	}

	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filter.Status = &status
		}
	}

	result, err := h.customerService.List(c.Request.Context(), filter)
	if err != nil {
		handleCustomerError(c, err)
		return
	}

	response.Success(c, CustomerListResponse{
		Customers: result.Customers,
		Total:     result.Total,
		Page:      page,
		Size:      pageSize,
	})
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid customer id")
		return
	}

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateCustomerInput{
		Name:    req.Name,
		Code:    req.Code,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
		Status:  req.Status,
	}

	customer, err := h.customerService.Update(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id, input)
	if err != nil {
		handleCustomerError(c, err)
		return
	}

	response.Success(c, customer)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid customer id")
		return
	}

	err = h.customerService.Delete(service.SetClientIPToContext(c.Request.Context(), middleware.GetClientIP(c)), id)
	if err != nil {
		handleCustomerError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleCustomerError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterCustomerRoutes(r *gin.RouterGroup, h *CustomerHandler) {
	customers := r.Group("/customers")
	{
		customers.GET("", h.List)
		customers.POST("", h.Create)
		customers.GET("/:id", h.GetByID)
		customers.PUT("/:id", h.Update)
		customers.DELETE("/:id", h.Delete)
	}
}
