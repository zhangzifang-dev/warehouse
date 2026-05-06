package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockOutboundOrderRepository struct {
	createFunc       func(ctx context.Context, order *model.OutboundOrder) error
	getByIDFunc      func(ctx context.Context, id int64) (*model.OutboundOrder, error)
	getByOrderNoFunc func(ctx context.Context, orderNo string) (*model.OutboundOrder, error)
	listFunc         func(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error)
	updateFunc       func(ctx context.Context, order *model.OutboundOrder) error
	deleteFunc       func(ctx context.Context, id int64) error
}

func (m *mockOutboundOrderRepository) Create(ctx context.Context, order *model.OutboundOrder) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order)
	}
	return errors.New("not implemented")
}

func (m *mockOutboundOrderRepository) GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOutboundOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.OutboundOrder, error) {
	if m.getByOrderNoFunc != nil {
		return m.getByOrderNoFunc(ctx, orderNo)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOutboundOrderRepository) List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, warehouseID, status)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockOutboundOrderRepository) Update(ctx context.Context, order *model.OutboundOrder) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, order)
	}
	return errors.New("not implemented")
}

func (m *mockOutboundOrderRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

type mockOutboundItemRepository struct {
	createFunc       func(ctx context.Context, item *model.OutboundItem) error
	listByOrderIDFunc func(ctx context.Context, orderID int64) ([]model.OutboundItem, error)
	updateFunc       func(ctx context.Context, item *model.OutboundItem) error
	deleteFunc       func(ctx context.Context, id int64) error
}

func (m *mockOutboundItemRepository) Create(ctx context.Context, item *model.OutboundItem) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, item)
	}
	return errors.New("not implemented")
}

func (m *mockOutboundItemRepository) ListByOrderID(ctx context.Context, orderID int64) ([]model.OutboundItem, error) {
	if m.listByOrderIDFunc != nil {
		return m.listByOrderIDFunc(ctx, orderID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockOutboundItemRepository) Update(ctx context.Context, item *model.OutboundItem) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, item)
	}
	return errors.New("not implemented")
}

func (m *mockOutboundItemRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestOutboundOrderService_Create_Success(t *testing.T) {
	createdOrder := &model.OutboundOrder{}
	mockOrderRepo := &mockOutboundOrderRepository{
		createFunc: func(ctx context.Context, order *model.OutboundOrder) error {
			order.ID = 1
			createdOrder = order
			return nil
		},
	}
	mockItemRepo := &mockOutboundItemRepository{}

	svc := NewOutboundOrderService(mockOrderRepo, mockItemRepo, nil, nil)
	input := &CreateOutboundOrderInput{
		OrderNo:       "SO-2024-001",
		WarehouseID:   1,
		TotalQuantity: 100,
	}

	order, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if order == nil {
		t.Fatal("expected order, got nil")
	}
	if createdOrder.OrderNo != "SO-2024-001" {
		t.Errorf("expected OrderNo 'SO-2024-001', got '%s'", createdOrder.OrderNo)
	}
}

func TestOutboundOrderService_Create_MissingWarehouseID(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{}
	mockItemRepo := &mockOutboundItemRepository{}

	svc := NewOutboundOrderService(mockOrderRepo, mockItemRepo, nil, nil)
	input := &CreateOutboundOrderInput{
		OrderNo:       "SO-2024-001",
		TotalQuantity: 100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing warehouse ID, got nil")
	}
}

func TestOutboundOrderService_GetByID_Success(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "SO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
			}, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	order, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if order == nil {
		t.Fatal("expected order, got nil")
	}
	if order.TotalQuantity != 100 {
		t.Errorf("expected total quantity 100, got %f", order.TotalQuantity)
	}
}

func TestOutboundOrderService_GetByID_NotFound(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent order, got nil")
	}
}

func TestOutboundOrderService_List_Success(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		listFunc: func(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error) {
			return []model.OutboundOrder{
				{BaseModel: model.BaseModel{ID: 1}, OrderNo: "SO-2024-001", WarehouseID: 1, TotalQuantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, OrderNo: "SO-2024-002", WarehouseID: 1, TotalQuantity: 200},
			}, 2, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	result, err := svc.List(context.Background(), 1, 10, 0, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result.Orders))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestOutboundOrderService_Update_Success(t *testing.T) {
	updatedOrder := &model.OutboundOrder{}
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "SO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
				Status:        0,
			}, nil
		},
		updateFunc: func(ctx context.Context, order *model.OutboundOrder) error {
			updatedOrder = order
			return nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)
	input := &UpdateOutboundOrderInput{
		TotalQuantity: floatPtr(200),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedOrder.TotalQuantity != 200 {
		t.Errorf("expected total quantity 200, got %f", updatedOrder.TotalQuantity)
	}
}

func TestOutboundOrderService_Delete_Success(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestOutboundOrderService_Delete_NotFound(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent order, got nil")
	}
}

func TestOutboundOrderService_Confirm_Success(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "SO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
				Status:        0,
			}, nil
		},
		updateFunc: func(ctx context.Context, order *model.OutboundOrder) error {
			return nil
		},
	}
	mockItemRepo := &mockOutboundItemRepository{
		listByOrderIDFunc: func(ctx context.Context, orderID int64) ([]model.OutboundItem, error) {
			return []model.OutboundItem{
				{OrderID: orderID, ProductID: 1, Quantity: 50, LocationID: 1},
				{OrderID: orderID, ProductID: 2, Quantity: 50, LocationID: 1},
			}, nil
		},
	}
	mockInventorySvc := &mockInventoryServiceForOutbound{
		checkStockFunc: func(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
			return &CheckStockResult{Available: true, CurrentStock: 100, Requested: input.Quantity}, nil
		},
		adjustQuantityFunc: func(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
			return &model.Inventory{}, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, mockItemRepo, mockInventorySvc)

	order, err := svc.Confirm(context.Background(), 1)

	if err != nil {
		t.Fatalf("Confirm failed: %v", err)
	}
	if order.Status != 1 {
		t.Errorf("expected status 1 (completed), got %d", order.Status)
	}
}

func TestOutboundOrderService_Confirm_InsufficientStock(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "SO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
				Status:        0,
			}, nil
		},
	}
	mockItemRepo := &mockOutboundItemRepository{
		listByOrderIDFunc: func(ctx context.Context, orderID int64) ([]model.OutboundItem, error) {
			return []model.OutboundItem{
				{OrderID: orderID, ProductID: 1, Quantity: 150, LocationID: 1},
			}, nil
		},
	}
	mockInventorySvc := &mockInventoryServiceForOutbound{
		checkStockFunc: func(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
			return &CheckStockResult{Available: false, CurrentStock: 100, Requested: 150}, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, mockItemRepo, mockInventorySvc)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for insufficient stock, got nil")
	}
}

func TestOutboundOrderService_Confirm_AlreadyCompleted(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel: model.BaseModel{ID: id},
				Status:    1,
			}, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for already completed order, got nil")
	}
}

func TestOutboundOrderService_Confirm_AlreadyCancelled(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel: model.BaseModel{ID: id},
				Status:    2,
			}, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for cancelled order, got nil")
	}
}

func TestOutboundOrderService_GetByOrderNo_Success(t *testing.T) {
	mockOrderRepo := &mockOutboundOrderRepository{
		getByOrderNoFunc: func(ctx context.Context, orderNo string) (*model.OutboundOrder, error) {
			return &model.OutboundOrder{
				BaseModel:     model.BaseModel{ID: 1},
				OrderNo:       orderNo,
				WarehouseID:   1,
				TotalQuantity: 100,
			}, nil
		},
	}

	svc := NewOutboundOrderService(mockOrderRepo, nil, nil, nil)

	order, err := svc.GetByOrderNo(context.Background(), "SO-2024-001")

	if err != nil {
		t.Fatalf("GetByOrderNo failed: %v", err)
	}
	if order.OrderNo != "SO-2024-001" {
		t.Errorf("expected OrderNo 'SO-2024-001', got '%s'", order.OrderNo)
	}
}

type mockInventoryServiceForOutbound struct {
	checkStockFunc    func(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error)
	adjustQuantityFunc func(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error)
}

func (m *mockInventoryServiceForOutbound) CheckStock(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
	if m.checkStockFunc != nil {
		return m.checkStockFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryServiceForOutbound) AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
	if m.adjustQuantityFunc != nil {
		return m.adjustQuantityFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}
