package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *model.Warehouse) error
	GetByID(ctx context.Context, id int64) (*model.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*model.Warehouse, error)
	List(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error)
	Update(ctx context.Context, warehouse *model.Warehouse) error
	Delete(ctx context.Context, id int64) error
}

type WarehouseService struct {
	warehouseRepo WarehouseRepository
}

func NewWarehouseService(warehouseRepo WarehouseRepository) *WarehouseService {
	return &WarehouseService{
		warehouseRepo: warehouseRepo,
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

	return warehouse, nil
}

func (s *WarehouseService) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "warehouse not found")
	}
	return warehouse, nil
}

func (s *WarehouseService) List(ctx context.Context, page, pageSize int) (*ListWarehousesResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	warehouses, total, err := s.warehouseRepo.List(ctx, page, pageSize)
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

	return warehouse, nil
}

func (s *WarehouseService) Delete(ctx context.Context, id int64) error {
	_, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "warehouse not found")
	}

	err = s.warehouseRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete warehouse")
	}

	return nil
}


