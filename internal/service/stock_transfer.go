package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type StockTransferRepository interface {
	Create(ctx context.Context, transfer *model.StockTransfer) error
	GetByID(ctx context.Context, id int64) (*model.StockTransfer, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*model.StockTransfer, error)
	List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) ([]model.StockTransfer, int, error)
	Update(ctx context.Context, transfer *model.StockTransfer) error
	Delete(ctx context.Context, id int64) error
}

type StockTransferItemRepository interface {
	Create(ctx context.Context, item *model.StockTransferItem) error
	ListByTransferID(ctx context.Context, transferID int64) ([]model.StockTransferItem, error)
	Update(ctx context.Context, item *model.StockTransferItem) error
	Delete(ctx context.Context, id int64) error
}

type InventoryServiceForTransfer interface {
	CheckStock(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error)
	AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error)
	GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error)
}

type StockTransferService struct {
	transferRepo StockTransferRepository
	itemRepo     StockTransferItemRepository
	inventorySvc InventoryServiceForTransfer
}

func NewStockTransferService(transferRepo StockTransferRepository, itemRepo StockTransferItemRepository, inventorySvc InventoryServiceForTransfer) *StockTransferService {
	return &StockTransferService{
		transferRepo: transferRepo,
		itemRepo:     itemRepo,
		inventorySvc: inventorySvc,
	}
}

type CreateStockTransferInput struct {
	OrderNo         string
	SourceWarehouseID int64
	TargetWarehouseID   int64
	TotalQty   float64
	Remark          string
}

type UpdateStockTransferInput struct {
	SourceWarehouseID *int64
	TargetWarehouseID   *int64
	TotalQty   *float64
	Status          *int
	Remark          *string
}

type ListStockTransfersResult struct {
	Transfers []model.StockTransfer
	Total     int
}

func (s *StockTransferService) Create(ctx context.Context, input *CreateStockTransferInput) (*model.StockTransfer, error) {
	if input.SourceWarehouseID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "from warehouse ID is required")
	}
	if input.TargetWarehouseID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "to warehouse ID is required")
	}
	if input.SourceWarehouseID == input.TargetWarehouseID {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "source and target warehouse cannot be the same")
	}

	transfer := &model.StockTransfer{
		OrderNo:           input.OrderNo,
		SourceWarehouseID: input.SourceWarehouseID,
		TargetWarehouseID: input.TargetWarehouseID,
		Status:            0,
	}

	err := s.transferRepo.Create(ctx, transfer)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create stock transfer")
	}

	return transfer, nil
}

func (s *StockTransferService) GetByID(ctx context.Context, id int64) (*model.StockTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "stock transfer not found")
	}
	return transfer, nil
}

func (s *StockTransferService) GetByOrderNo(ctx context.Context, orderNo string) (*model.StockTransfer, error) {
	transfer, err := s.transferRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "stock transfer not found")
	}
	return transfer, nil
}

func (s *StockTransferService) List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) (*ListStockTransfersResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	transfers, total, err := s.transferRepo.List(ctx, page, pageSize, fromWarehouseID, toWarehouseID, status)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list stock transfers")
	}

	return &ListStockTransfersResult{
		Transfers: transfers,
		Total:     total,
	}, nil
}

func (s *StockTransferService) Update(ctx context.Context, id int64, input *UpdateStockTransferInput) (*model.StockTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "stock transfer not found")
	}

	if input.SourceWarehouseID != nil {
		transfer.SourceWarehouseID = *input.SourceWarehouseID
	}
	if input.TargetWarehouseID != nil {
		transfer.TargetWarehouseID = *input.TargetWarehouseID
	}
	if input.Status != nil {
		transfer.Status = *input.Status
	}

	if transfer.SourceWarehouseID == transfer.TargetWarehouseID {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "source and target warehouse cannot be the same")
	}

	err = s.transferRepo.Update(ctx, transfer)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update stock transfer")
	}

	return transfer, nil
}

func (s *StockTransferService) Delete(ctx context.Context, id int64) error {
	_, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "stock transfer not found")
	}

	err = s.transferRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete stock transfer")
	}

	return nil
}

func (s *StockTransferService) Confirm(ctx context.Context, id int64) (*model.StockTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "stock transfer not found")
	}

	if transfer.Status == 1 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "transfer already completed")
	}

	if transfer.Status == 2 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "transfer already cancelled")
	}

	if s.itemRepo != nil && s.inventorySvc != nil {
		items, err := s.itemRepo.ListByTransferID(ctx, id)
		if err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get transfer items")
		}

		for _, item := range items {
			checkResult, err := s.inventorySvc.CheckStock(ctx, &CheckStockInput{
				WarehouseID: transfer.SourceWarehouseID,
				ProductID:   item.ProductID,
				BatchNo:     item.BatchNo,
				Quantity:    item.Quantity,
			})
			if err != nil {
				return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to check inventory")
			}

			if !checkResult.Available {
				return nil, apperrors.NewAppError(apperrors.CodeInsufficientStock, "insufficient stock in source warehouse")
			}
		}

		for _, item := range items {
			sourceInventory, err := s.inventorySvc.GetByWarehouseAndProduct(ctx, transfer.SourceWarehouseID, item.ProductID, item.BatchNo)
			if err != nil {
				return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get source inventory")
			}

			_, err = s.inventorySvc.AdjustQuantity(ctx, &AdjustQuantityInput{
				InventoryID: sourceInventory.ID,
				Quantity:    -item.Quantity,
			})
			if err != nil {
				return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to decrease stock in source warehouse")
			}
		}

		for _, item := range items {
			targetInventory, err := s.inventorySvc.GetByWarehouseAndProduct(ctx, transfer.TargetWarehouseID, item.ProductID, item.BatchNo)
			if err != nil {
				targetInventory = &model.Inventory{
					WarehouseID: transfer.TargetWarehouseID,
					ProductID:   item.ProductID,
					LocationID:  item.LocationID,
					Quantity:    0,
					BatchNo:     item.BatchNo,
				}
			}

			_, err = s.inventorySvc.AdjustQuantity(ctx, &AdjustQuantityInput{
				InventoryID: targetInventory.ID,
				Quantity:    item.Quantity,
			})
			if err != nil {
				return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to increase stock in target warehouse")
			}
		}
	}

	transfer.Status = 1

	err = s.transferRepo.Update(ctx, transfer)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update transfer status")
	}

	return transfer, nil
}
