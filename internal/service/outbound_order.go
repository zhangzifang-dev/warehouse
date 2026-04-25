package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type OutboundOrderRepository interface {
	Create(ctx context.Context, order *model.OutboundOrder) error
	GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*model.OutboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error)
	Update(ctx context.Context, order *model.OutboundOrder) error
	Delete(ctx context.Context, id int64) error
}

type OutboundItemRepository interface {
	Create(ctx context.Context, item *model.OutboundItem) error
	ListByOrderID(ctx context.Context, orderID int64) ([]model.OutboundItem, error)
	Update(ctx context.Context, item *model.OutboundItem) error
	Delete(ctx context.Context, id int64) error
}

type InventoryServiceForOutbound interface {
	CheckStock(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error)
	AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error)
}

type OutboundOrderService struct {
	orderRepo     OutboundOrderRepository
	itemRepo      OutboundItemRepository
	inventorySvc   InventoryServiceForOutbound
}

func NewOutboundOrderService(orderRepo OutboundOrderRepository, itemRepo OutboundItemRepository, inventorySvc InventoryServiceForOutbound) *OutboundOrderService {
	return &OutboundOrderService{
		orderRepo:    orderRepo,
		itemRepo:     itemRepo,
		inventorySvc:  inventorySvc,
	}
}

type CreateOutboundOrderInput struct {
	OrderNo       string
	CustomerID    int64
	WarehouseID   int64
	TotalQuantity float64
	Remark        string
}

type UpdateOutboundOrderInput struct {
	CustomerID    *int64
	WarehouseID   *int64
	TotalQuantity *float64
	Status        *int
	Remark        *string
}

type ListOutboundOrdersResult struct {
	Orders []model.OutboundOrder
	Total  int
}

func (s *OutboundOrderService) Create(ctx context.Context, input *CreateOutboundOrderInput) (*model.OutboundOrder, error) {
	if input.WarehouseID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse ID is required")
	}

	order := &model.OutboundOrder{
		OrderNo:       input.OrderNo,
		CustomerID:    input.CustomerID,
		WarehouseID:   input.WarehouseID,
		TotalQuantity: input.TotalQuantity,
		Status:        0,
		Remark:        input.Remark,
	}

	err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create outbound order")
	}

	return order, nil
}

func (s *OutboundOrderService) GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "outbound order not found")
	}
	return order, nil
}

func (s *OutboundOrderService) GetByOrderNo(ctx context.Context, orderNo string) (*model.OutboundOrder, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "outbound order not found")
	}
	return order, nil
}

func (s *OutboundOrderService) List(ctx context.Context, page, pageSize int, warehouseID, status int) (*ListOutboundOrdersResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	orders, total, err := s.orderRepo.List(ctx, page, pageSize, warehouseID, status)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list outbound orders")
	}

	return &ListOutboundOrdersResult{
		Orders: orders,
		Total:  total,
	}, nil
}

func (s *OutboundOrderService) Update(ctx context.Context, id int64, input *UpdateOutboundOrderInput) (*model.OutboundOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "outbound order not found")
	}

	if input.CustomerID != nil {
		order.CustomerID = *input.CustomerID
	}
	if input.WarehouseID != nil {
		order.WarehouseID = *input.WarehouseID
	}
	if input.TotalQuantity != nil {
		order.TotalQuantity = *input.TotalQuantity
	}
	if input.Status != nil {
		order.Status = *input.Status
	}
	if input.Remark != nil {
		order.Remark = *input.Remark
	}

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update outbound order")
	}

	return order, nil
}

func (s *OutboundOrderService) Delete(ctx context.Context, id int64) error {
	_, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "outbound order not found")
	}

	err = s.orderRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete outbound order")
	}

	return nil
}

func (s *OutboundOrderService) Confirm(ctx context.Context, id int64) (*model.OutboundOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "outbound order not found")
	}

	if order.Status == 1 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "order already completed")
	}

	if order.Status == 2 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "order already cancelled")
	}

	if s.itemRepo != nil && s.inventorySvc != nil {
		items, err := s.itemRepo.ListByOrderID(ctx, id)
		if err != nil {
			return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get order items")
		}

		for _, item := range items {
			checkResult, err := s.inventorySvc.CheckStock(ctx, &CheckStockInput{
				WarehouseID: order.WarehouseID,
				ProductID:   item.ProductID,
				BatchNo:     item.BatchNo,
				Quantity:    item.Quantity,
			})
			if err != nil {
				return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to check inventory")
			}

			if !checkResult.Available {
				return nil, apperrors.NewAppError(apperrors.CodeInsufficientStock, "insufficient stock")
			}
		}

		for _, item := range items {
			_, err := s.inventorySvc.AdjustQuantity(ctx, &AdjustQuantityInput{
				InventoryID: 0,
				Quantity:    -item.Quantity,
			})
			if err != nil {
				return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update inventory")
			}
		}
	}

	order.Status = 1

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update order status")
	}

	return order, nil
}
