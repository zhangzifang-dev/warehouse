package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type OutboundOrderRepository interface {
	Create(ctx context.Context, order *model.OutboundOrder) error
	GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*model.OutboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error)
	ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) ([]model.OutboundOrder, int, error)
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
	orderRepo    OutboundOrderRepository
	itemRepo     OutboundItemRepository
	inventorySvc InventoryServiceForOutbound
	auditLogger  AuditLogger
}

func NewOutboundOrderService(orderRepo OutboundOrderRepository, itemRepo OutboundItemRepository, inventorySvc InventoryServiceForOutbound, auditLogger AuditLogger) *OutboundOrderService {
	return &OutboundOrderService{
		orderRepo:    orderRepo,
		itemRepo:     itemRepo,
		inventorySvc: inventorySvc,
		auditLogger:  auditLogger,
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

	if s.auditLogger != nil {
		jsonBytes, _ := json.Marshal(order)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "outbound_orders",
			RecordID:   order.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(jsonBytes)},
			OperatedBy: order.CreatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
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

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(order)
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

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(order)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "outbound_orders",
			RecordID:   order.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: order.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return order, nil
}

func (s *OutboundOrderService) Delete(ctx context.Context, id int64) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "outbound order not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(order)
	}

	err = s.orderRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete outbound order")
	}

	if s.auditLogger != nil {
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "outbound_orders",
			RecordID:   order.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: order.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
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

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(order)
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

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(order)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "outbound_orders",
			RecordID:   order.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: order.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return order, nil
}

func (s *OutboundOrderService) ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) (*ListOutboundOrdersResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	orders, total, err := s.orderRepo.ListWithFilter(ctx, filter)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list outbound orders")
	}

	return &ListOutboundOrdersResult{
		Orders: orders,
		Total:  total,
	}, nil
}
