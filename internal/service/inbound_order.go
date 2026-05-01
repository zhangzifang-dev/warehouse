package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type InboundOrderRepository interface {
	Create(ctx context.Context, order *model.InboundOrder) error
	GetByID(ctx context.Context, id int64) (*model.InboundOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*model.InboundOrder, error)
	List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.InboundOrder, int, error)
	Update(ctx context.Context, order *model.InboundOrder) error
	Delete(ctx context.Context, id int64) error
}

type InboundItemRepository interface {
	Create(ctx context.Context, item *model.InboundItem) error
	ListByOrderID(ctx context.Context, orderID int64) ([]model.InboundItem, error)
	Update(ctx context.Context, item *model.InboundItem) error
	Delete(ctx context.Context, id int64) error
}

type InventoryServiceForInbound interface {
	AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error)
}

type InboundOrderService struct {
	orderRepo    InboundOrderRepository
	itemRepo     InboundItemRepository
	inventorySvc InventoryServiceForInbound
	auditLogger  AuditLogger
}

func NewInboundOrderService(orderRepo InboundOrderRepository, itemRepo InboundItemRepository, inventorySvc InventoryServiceForInbound, auditLogger AuditLogger) *InboundOrderService {
	return &InboundOrderService{
		orderRepo:    orderRepo,
		itemRepo:     itemRepo,
		inventorySvc: inventorySvc,
		auditLogger:  auditLogger,
	}
}

type CreateInboundOrderInput struct {
	OrderNo       string
	SupplierID    int64
	WarehouseID   int64
	TotalQuantity float64
	Remark        string
}

type UpdateInboundOrderInput struct {
	SupplierID    *int64
	WarehouseID   *int64
	TotalQuantity *float64
	Status        *int
	Remark        *string
}

type ListInboundOrdersResult struct {
	Orders []model.InboundOrder
	Total  int
}

func (s *InboundOrderService) Create(ctx context.Context, input *CreateInboundOrderInput) (*model.InboundOrder, error) {
	if input.WarehouseID <= 0 {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "warehouse ID is required")
	}

	order := &model.InboundOrder{
		OrderNo:       input.OrderNo,
		SupplierID:    input.SupplierID,
		WarehouseID:   input.WarehouseID,
		TotalQuantity: input.TotalQuantity,
		Status:        0,
		Remark:        input.Remark,
	}

	err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create inbound order")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(order)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "inbound_orders",
			RecordID:   order.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: order.CreatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return order, nil
}

func (s *InboundOrderService) GetByID(ctx context.Context, id int64) (*model.InboundOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inbound order not found")
	}
	return order, nil
}

func (s *InboundOrderService) GetByOrderNo(ctx context.Context, orderNo string) (*model.InboundOrder, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inbound order not found")
	}
	return order, nil
}

func (s *InboundOrderService) List(ctx context.Context, page, pageSize int, warehouseID, status int) (*ListInboundOrdersResult, error) {
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
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list inbound orders")
	}

	return &ListInboundOrdersResult{
		Orders: orders,
		Total:  total,
	}, nil
}

func (s *InboundOrderService) Update(ctx context.Context, id int64, input *UpdateInboundOrderInput) (*model.InboundOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inbound order not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(order)
	}

	if input.SupplierID != nil {
		order.SupplierID = *input.SupplierID
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
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update inbound order")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(order)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "inbound_orders",
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

func (s *InboundOrderService) Delete(ctx context.Context, id int64) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "inbound order not found")
	}

	err = s.orderRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete inbound order")
	}

	if s.auditLogger != nil {
		oldValue, _ := json.Marshal(order)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "inbound_orders",
			RecordID:   order.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: order.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return nil
}

func (s *InboundOrderService) Confirm(ctx context.Context, id int64) (*model.InboundOrder, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "inbound order not found")
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
			_, err := s.inventorySvc.AdjustQuantity(ctx, &AdjustQuantityInput{
				InventoryID: 0,
				Quantity:    item.Quantity,
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
			TableName:  "inbound_orders",
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
