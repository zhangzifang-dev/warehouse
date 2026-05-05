package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	"warehouse/internal/repository"
	apperrors "warehouse/internal/pkg/errors"
)

type WarehouseQueryFilter struct {
	Name     string
	Page     int
	PageSize int
}

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *model.Warehouse) error
	GetByID(ctx context.Context, id int64) (*model.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*model.Warehouse, error)
	List(ctx context.Context, filter *repository.WarehouseQueryFilter) ([]model.Warehouse, int, error)
	Update(ctx context.Context, warehouse *model.Warehouse) error
	Delete(ctx context.Context, id int64) error
}

type AuditLogger interface {
	Log(ctx context.Context, input *CreateAuditLogInput) error
}

type WarehouseService struct {
	warehouseRepo WarehouseRepository
	auditLogger   AuditLogger
}

func NewWarehouseService(warehouseRepo WarehouseRepository, auditLogger AuditLogger) *WarehouseService {
	return &WarehouseService{
		warehouseRepo: warehouseRepo,
		auditLogger:   auditLogger,
	}
}

type CreateWarehouseInput struct {
	Name    string
	Code    string
	Address string
	Contact string
	Phone   string
	Status  int
}

type UpdateWarehouseInput struct {
	Name    string
	Address string
	Contact string
	Phone   string
	Status  *int
}

type ListWarehousesResult struct {
	Warehouses []model.Warehouse
	Total      int
}

func (s *WarehouseService) Create(ctx context.Context, input *CreateWarehouseInput) (*model.Warehouse, error) {
	existing, err := s.warehouseRepo.GetByCode(ctx, input.Code)
	if err == nil && existing != nil {
		return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "warehouse code already exists")
	}

	warehouse := &model.Warehouse{
		Name:    input.Name,
		Code:    input.Code,
		Address: input.Address,
		Contact: input.Contact,
		Phone:   input.Phone,
		Status:  input.Status,
	}

	if warehouse.Status == 0 {
		warehouse.Status = model.WarehouseStatusActive
	}

	err = s.warehouseRepo.Create(ctx, warehouse)
	if err != nil {
		if isDuplicateEntry(err) {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "warehouse code already exists")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create warehouse")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(warehouse)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "warehouses",
			RecordID:   warehouse.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: warehouse.CreatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return warehouse, nil
}

func (s *WarehouseService) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "warehouse not found")
	}
	return warehouse, nil
}

func (s *WarehouseService) List(ctx context.Context, filter *WarehouseQueryFilter) (*ListWarehousesResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	repoFilter := &repository.WarehouseQueryFilter{
		Name:     filter.Name,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}

	warehouses, total, err := s.warehouseRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list warehouses")
	}

	return &ListWarehousesResult{
		Warehouses: warehouses,
		Total:      total,
	}, nil
}

func (s *WarehouseService) Update(ctx context.Context, id int64, input *UpdateWarehouseInput) (*model.Warehouse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "warehouse not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(warehouse)
	}

	if input.Name != "" {
		warehouse.Name = input.Name
	}
	if input.Address != "" {
		warehouse.Address = input.Address
	}
	if input.Contact != "" {
		warehouse.Contact = input.Contact
	}
	if input.Phone != "" {
		warehouse.Phone = input.Phone
	}
	if input.Status != nil {
		warehouse.Status = *input.Status
	}

	err = s.warehouseRepo.Update(ctx, warehouse)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update warehouse")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(warehouse)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "warehouses",
			RecordID:   warehouse.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: warehouse.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return warehouse, nil
}

func (s *WarehouseService) Delete(ctx context.Context, id int64) error {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "warehouse not found")
	}

	err = s.warehouseRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete warehouse")
	}

	if s.auditLogger != nil {
		oldValue, _ := json.Marshal(warehouse)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "warehouses",
			RecordID:   warehouse.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: warehouse.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return nil
}


