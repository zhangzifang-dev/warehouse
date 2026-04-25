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

type CreateProductRequest struct {
	SKU           string   `json:"sku" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	CategoryID    *int64   `json:"category_id"`
	Specification *string  `json:"specification"`
	Unit          *string  `json:"unit"`
	Barcode       *string  `json:"barcode"`
	Price         *float64 `json:"price"`
	Description   *string  `json:"description"`
	Status        *int     `json:"status"`
}

type UpdateProductRequest struct {
	SKU           *string  `json:"sku"`
	Name          *string  `json:"name"`
	CategoryID    *int64   `json:"category_id"`
	Specification *string  `json:"specification"`
	Unit          *string  `json:"unit"`
	Barcode       *string  `json:"barcode"`
	Price         *float64 `json:"price"`
	Description   *string  `json:"description"`
	Status        *int     `json:"status"`
}

type ProductListResponse struct {
	Products []model.Product `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	Size     int             `json:"size"`
}

type ProductService interface {
	Create(ctx context.Context, input *service.CreateProductInput) (*model.Product, error)
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) (*service.ListProductsResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateProductInput) (*model.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(productService ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateProductInput{
		SKU:  req.SKU,
		Name: req.Name,
	}

	if req.CategoryID != nil {
		input.CategoryID = *req.CategoryID
	}
	if req.Specification != nil {
		input.Specification = *req.Specification
	}
	if req.Unit != nil {
		input.Unit = *req.Unit
	}
	if req.Barcode != nil {
		input.Barcode = *req.Barcode
	}
	if req.Price != nil {
		input.Price = *req.Price
	}
	if req.Description != nil {
		input.Description = *req.Description
	}
	if req.Status != nil {
		input.Status = *req.Status
	}

	product, err := h.productService.Create(c.Request.Context(), input)
	if err != nil {
		handleProductError(c, err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid product id")
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleProductError(c, err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)
	keyword := c.Query("keyword")

	result, err := h.productService.List(c.Request.Context(), page, pageSize, categoryID, keyword)
	if err != nil {
		handleProductError(c, err)
		return
	}

	response.Success(c, ProductListResponse{
		Products: result.Products,
		Total:    result.Total,
		Page:     page,
		Size:     pageSize,
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid product id")
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateProductInput{
		SKU:           req.SKU,
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Specification: req.Specification,
		Unit:          req.Unit,
		Barcode:       req.Barcode,
		Price:         req.Price,
		Description:   req.Description,
		Status:        req.Status,
	}

	product, err := h.productService.Update(c.Request.Context(), id, input)
	if err != nil {
		handleProductError(c, err)
		return
	}

	response.Success(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid product id")
		return
	}

	err = h.productService.Delete(c.Request.Context(), id)
	if err != nil {
		handleProductError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleProductError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterProductRoutes(r *gin.RouterGroup, h *ProductHandler) {
	products := r.Group("/products")
	{
		products.GET("", h.List)
		products.POST("", h.Create)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}
