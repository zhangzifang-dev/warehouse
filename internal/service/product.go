package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	GetBySKU(ctx context.Context, sku string) (*model.Product, error)
	List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int64) error
}

type ProductService struct {
	productRepo ProductRepository
	auditLogger AuditLogger
}

func NewProductService(productRepo ProductRepository, auditLogger AuditLogger) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		auditLogger: auditLogger,
	}
}

type CreateProductInput struct {
	SKU           string
	Name          string
	CategoryID    int64
	Specification string
	Unit          string
	Barcode       string
	Price         float64
	Description   string
	Status        int
}

type UpdateProductInput struct {
	SKU           *string
	Name          *string
	CategoryID    *int64
	Specification *string
	Unit          *string
	Barcode       *string
	Price         *float64
	Description   *string
	Status        *int
}

type ListProductsResult struct {
	Products []model.Product
	Total    int
}

func (s *ProductService) Create(ctx context.Context, input *CreateProductInput) (*model.Product, error) {
	if input.SKU == "" {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "product SKU is required")
	}
	if input.Name == "" {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "product name is required")
	}

	existing, err := s.productRepo.GetBySKU(ctx, input.SKU)
	if err == nil && existing != nil {
		return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "SKU already exists")
	}

	product := &model.Product{
		SKU:           input.SKU,
		Name:          input.Name,
		CategoryID:    input.CategoryID,
		Specification: input.Specification,
		Unit:          input.Unit,
		Barcode:       input.Barcode,
		Price:         input.Price,
		Description:   input.Description,
		Status:        input.Status,
	}

	if product.Status == 0 {
		product.Status = model.ProductStatusActive
	}

	err = s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create product")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(product)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "products",
			RecordID:   product.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: product.CreatedBy,
		})
	}

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "product not found")
	}
	return product, nil
}

func (s *ProductService) List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) (*ListProductsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	products, total, err := s.productRepo.List(ctx, page, pageSize, categoryID, keyword)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list products")
	}

	return &ListProductsResult{
		Products: products,
		Total:    total,
	}, nil
}

func (s *ProductService) Update(ctx context.Context, id int64, input *UpdateProductInput) (*model.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "product not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(product)
	}

	if input.SKU != nil && *input.SKU != product.SKU {
		existing, err := s.productRepo.GetBySKU(ctx, *input.SKU)
		if err == nil && existing != nil {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "SKU already exists")
		}
		product.SKU = *input.SKU
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	if input.Specification != nil {
		product.Specification = *input.Specification
	}
	if input.Unit != nil {
		product.Unit = *input.Unit
	}
	if input.Barcode != nil {
		product.Barcode = *input.Barcode
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Status != nil {
		product.Status = *input.Status
	}

	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update product")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(product)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "products",
			RecordID:   product.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: product.UpdatedBy,
		})
	}

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "product not found")
	}

	err = s.productRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete product")
	}

	if s.auditLogger != nil {
		oldValue, _ := json.Marshal(product)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "products",
			RecordID:   product.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: product.UpdatedBy,
		})
	}

	return nil
}
