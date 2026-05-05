package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/repository"
)

type SupplierQueryFilter struct {
	Code     string
	Name     string
	Contact  string
	Phone    string
	Status   *int
	Page     int
	PageSize int
}

type SupplierRepository interface {
	Create(ctx context.Context, supplier *model.Supplier) error
	GetByID(ctx context.Context, id int64) (*model.Supplier, error)
	GetByCode(ctx context.Context, code string) (*model.Supplier, error)
	List(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error)
	Update(ctx context.Context, supplier *model.Supplier) error
	Delete(ctx context.Context, id int64) error
}

type SupplierService struct {
	supplierRepo SupplierRepository
	auditLogger  AuditLogger
}

func NewSupplierService(supplierRepo SupplierRepository, auditLogger AuditLogger) *SupplierService {
	return &SupplierService{
		supplierRepo: supplierRepo,
		auditLogger:  auditLogger,
	}
}

type CreateSupplierInput struct {
	Name    string
	Code    string
	Contact string
	Phone   string
	Email   string
	Address string
	Status  int
}

type UpdateSupplierInput struct {
	Name    *string
	Code    *string
	Contact *string
	Phone   *string
	Email   *string
	Address *string
	Status  *int
}

type ListSuppliersResult struct {
	Suppliers []model.Supplier
	Total     int
}

func (s *SupplierService) Create(ctx context.Context, input *CreateSupplierInput) (*model.Supplier, error) {
	if input.Name == "" {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "supplier name is required")
	}

	if input.Code != "" {
		existing, err := s.supplierRepo.GetByCode(ctx, input.Code)
		if err == nil && existing != nil {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "supplier code already exists")
		}
	}

	supplier := &model.Supplier{
		Name:    input.Name,
		Code:    input.Code,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Address: input.Address,
		Status:  input.Status,
	}

	if supplier.Status == 0 {
		supplier.Status = model.SupplierStatusActive
	}

	err := s.supplierRepo.Create(ctx, supplier)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create supplier")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(supplier)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "suppliers",
			RecordID:   supplier.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: supplier.CreatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return supplier, nil
}

func (s *SupplierService) GetByID(ctx context.Context, id int64) (*model.Supplier, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "supplier not found")
	}
	return supplier, nil
}

func (s *SupplierService) List(ctx context.Context, filter *SupplierQueryFilter) (*ListSuppliersResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	repoFilter := &repository.SupplierQueryFilter{
		Code:     filter.Code,
		Name:     filter.Name,
		Contact:  filter.Contact,
		Phone:    filter.Phone,
		Status:   filter.Status,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}

	suppliers, total, err := s.supplierRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list suppliers")
	}

	return &ListSuppliersResult{
		Suppliers: suppliers,
		Total:     total,
	}, nil
}

func (s *SupplierService) Update(ctx context.Context, id int64, input *UpdateSupplierInput) (*model.Supplier, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "supplier not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(supplier)
	}

	if input.Code != nil && *input.Code != supplier.Code {
		if *input.Code != "" {
			existing, err := s.supplierRepo.GetByCode(ctx, *input.Code)
			if err == nil && existing != nil {
				return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "supplier code already exists")
			}
		}
		supplier.Code = *input.Code
	}

	if input.Name != nil {
		supplier.Name = *input.Name
	}
	if input.Contact != nil {
		supplier.Contact = *input.Contact
	}
	if input.Phone != nil {
		supplier.Phone = *input.Phone
	}
	if input.Email != nil {
		supplier.Email = *input.Email
	}
	if input.Address != nil {
		supplier.Address = *input.Address
	}
	if input.Status != nil {
		supplier.Status = *input.Status
	}

	err = s.supplierRepo.Update(ctx, supplier)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update supplier")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(supplier)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "suppliers",
			RecordID:   supplier.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: supplier.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return supplier, nil
}

func (s *SupplierService) Delete(ctx context.Context, id int64) error {
	supplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "supplier not found")
	}

	err = s.supplierRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete supplier")
	}

	if s.auditLogger != nil {
		oldValue, _ := json.Marshal(supplier)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "suppliers",
			RecordID:   supplier.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: supplier.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return nil
}
