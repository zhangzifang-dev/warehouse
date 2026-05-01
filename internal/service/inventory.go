package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory *model.Inventory) error
	GetByID(ctx context.Context, id int64) (*model.Inventory, error)
	List(ctx context.Context, page, pageSize int, warehouseID, productID int64) ([]model.Inventory, int, error)
	Update(ctx context.Context, inventory *model.Inventory) error
	Delete(ctx context.Context, id int64) error
	GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error)
	UpdateQuantity(ctx context.Context, id int64, quantity float64) error
}

type InventoryService struct {
	inventoryRepo InventoryRepository
	auditLogger   AuditLogger
}

func NewInventoryService(inventoryRepo InventoryRepository, auditLogger AuditLogger) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		auditLogger:   auditLogger,
	}
}

type CreateInventoryInput struct {
	WarehouseID int64
	ProductID   int64
	LocationID  int64
	Quantity    float64
	BatchNo     string
}

type UpdateInventoryInput struct {
	WarehouseID *int64
	ProductID   *int64
	LocationID  *int64
	Quantity    *float64
	BatchNo     *string
}

type ListInventoriesResult struct {
	Inventories []model.Inventory
	Total       int
}

type AdjustQuantityInput struct {
	InventoryID int64
	Quantity    float64
}

type CheckStockInput struct {
	WarehouseID int64
	ProductID   int64
	BatchNo     string
	Quantity    float64
}

type CheckStockResult struct {
	Available     bool
	CurrentStock  float64
	Requested     float64
}

func (s *InventoryService) Create(ctx context.Context, input *CreateInventoryInput) (*model.Inventory, error) {
	if input.WarehouseID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse ID is required")
	}
	if input.ProductID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "product ID is required")
	}

	inventory := &model.Inventory{
		WarehouseID: input.WarehouseID,
		ProductID:   input.ProductID,
		LocationID:  input.LocationID,
		Quantity:    input.Quantity,
		BatchNo:     input.BatchNo,
	}

	err := s.inventoryRepo.Create(ctx, inventory)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create inventory")
	}

	return inventory, nil
}

func (s *InventoryService) GetByID(ctx context.Context, id int64) (*model.Inventory, error) {
	inventory, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inventory not found")
	}
	return inventory, nil
}

func (s *InventoryService) List(ctx context.Context, page, pageSize int, warehouseID, productID int64) (*ListInventoriesResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	inventories, total, err := s.inventoryRepo.List(ctx, page, pageSize, warehouseID, productID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list inventories")
	}

	return &ListInventoriesResult{
		Inventories: inventories,
		Total:       total,
	}, nil
}

func (s *InventoryService) Update(ctx context.Context, id int64, input *UpdateInventoryInput) (*model.Inventory, error) {
	inventory, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inventory not found")
	}

	if input.WarehouseID != nil {
		inventory.WarehouseID = *input.WarehouseID
	}
	if input.ProductID != nil {
		inventory.ProductID = *input.ProductID
	}
	if input.LocationID != nil {
		inventory.LocationID = *input.LocationID
	}
	if input.Quantity != nil {
		inventory.Quantity = *input.Quantity
	}
	if input.BatchNo != nil {
		inventory.BatchNo = *input.BatchNo
	}

	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update inventory")
	}

	return inventory, nil
}

func (s *InventoryService) GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
	inventory, err := s.inventoryRepo.GetByWarehouseAndProduct(ctx, warehouseID, productID, batchNo)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *InventoryService) Delete(ctx context.Context, id int64) error {
	_, err := s.inventoryRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "inventory not found")
	}

	err = s.inventoryRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete inventory")
	}

	return nil
}

func (s *InventoryService) AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
	if input.InventoryID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "inventory ID is required")
	}

	inventory, err := s.inventoryRepo.GetByID(ctx, input.InventoryID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inventory not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(inventory)
	}

	newQuantity := inventory.Quantity + input.Quantity

	if newQuantity < 0 {
		return nil, apperrors.NewAppError(apperrors.CodeInsufficientStock, "insufficient stock")
	}

	inventory.Quantity = newQuantity

	err = s.inventoryRepo.Update(ctx, inventory)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to adjust quantity")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(inventory)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "inventory",
			RecordID:   inventory.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: inventory.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return inventory, nil
}

func (s *InventoryService) CheckStock(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
	if input.WarehouseID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse ID is required")
	}
	if input.ProductID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "product ID is required")
	}
	if input.Quantity <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "quantity must be positive")
	}

	inventory, err := s.inventoryRepo.GetByWarehouseAndProduct(ctx, input.WarehouseID, input.ProductID, input.BatchNo)
	if err != nil {
		return &CheckStockResult{
			Available:    false,
			CurrentStock: 0,
			Requested:    input.Quantity,
		}, nil
	}

	available := inventory.Quantity >= input.Quantity

	return &CheckStockResult{
		Available:    available,
		CurrentStock: inventory.Quantity,
		Requested:    input.Quantity,
	}, nil
}
