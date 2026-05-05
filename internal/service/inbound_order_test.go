package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockInboundOrderRepository struct {
	createFunc      func(ctx context.Context, order *model.InboundOrder) error
	getByIDFunc     func(ctx context.Context, id int64) (*model.InboundOrder, error)
	getByOrderNoFunc func(ctx context.Context, orderNo string) (*model.InboundOrder, error)
	listFunc        func(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.InboundOrder, int, error)
	updateFunc      func(ctx context.Context, order *model.InboundOrder) error
	deleteFunc      func(ctx context.Context, id int64) error
}

func (m *mockInboundOrderRepository) Create(ctx context.Context, order *model.InboundOrder) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order)
	}
	return errors.New("not implemented")
}

func (m *mockInboundOrderRepository) GetByID(ctx context.Context, id int64) (*model.InboundOrder, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.InboundOrder, error) {
	if m.getByOrderNoFunc != nil {
		return m.getByOrderNoFunc(ctx, orderNo)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundOrderRepository) List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.InboundOrder, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, warehouseID, status)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockInboundOrderRepository) Update(ctx context.Context, order *model.InboundOrder) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, order)
	}
	return errors.New("not implemented")
}

func (m *mockInboundOrderRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

type mockInboundItemRepository struct {
	createFunc       func(ctx context.Context, item *model.InboundItem) error
	listByOrderIDFunc func(ctx context.Context, orderID int64) ([]model.InboundItem, error)
	updateFunc       func(ctx context.Context, item *model.InboundItem) error
	deleteFunc       func(ctx context.Context, id int64) error
}

func (m *mockInboundItemRepository) Create(ctx context.Context, item *model.InboundItem) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, item)
	}
	return errors.New("not implemented")
}

func (m *mockInboundItemRepository) ListByOrderID(ctx context.Context, orderID int64) ([]model.InboundItem, error) {
	if m.listByOrderIDFunc != nil {
		return m.listByOrderIDFunc(ctx, orderID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInboundItemRepository) Update(ctx context.Context, item *model.InboundItem) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, item)
	}
	return errors.New("not implemented")
}

func (m *mockInboundItemRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestInboundOrderService_Create_Success(t *testing.T) {
	createdOrder := &model.InboundOrder{}
	mockOrderRepo := &mockInboundOrderRepository{
		createFunc: func(ctx context.Context, order *model.InboundOrder) error {
			order.ID = 1
			createdOrder = order
			return nil
		},
	}
	mockItemRepo := &mockInboundItemRepository{}

	svc := NewInboundOrderService(mockOrderRepo, mockItemRepo, nil, nil)
	input := &CreateInboundOrderInput{
		OrderNo:       "PO-2024-001",
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
	if createdOrder.OrderNo != "PO-2024-001" {
		t.Errorf("expected OrderNo 'PO-2024-001', got '%s'", createdOrder.OrderNo)
	}
}

func TestInboundOrderService_Create_MissingWarehouseID(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{}
	mockItemRepo := &mockInboundItemRepository{}

	svc := NewInboundOrderService(mockOrderRepo, mockItemRepo, nil, nil)
	input := &CreateInboundOrderInput{
		OrderNo:       "PO-2024-001",
		TotalQuantity: 100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing warehouse ID, got nil")
	}
}

func TestInboundOrderService_GetByID_Success(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return &model.InboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "PO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
			}, nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

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

func TestInboundOrderService_GetByID_NotFound(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent order, got nil")
	}
}

func TestInboundOrderService_List_Success(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		listFunc: func(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.InboundOrder, int, error) {
			return []model.InboundOrder{
				{BaseModel: model.BaseModel{ID: 1}, OrderNo: "PO-2024-001", WarehouseID: 1, TotalQuantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, OrderNo: "PO-2024-002", WarehouseID: 1, TotalQuantity: 200},
			}, 2, nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

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

func TestInboundOrderService_Update_Success(t *testing.T) {
	updatedOrder := &model.InboundOrder{}
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return &model.InboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "PO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
				Status:        0,
			}, nil
		},
		updateFunc: func(ctx context.Context, order *model.InboundOrder) error {
			updatedOrder = order
			return nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)
	input := &UpdateInboundOrderInput{
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

func TestInboundOrderService_Delete_Success(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return &model.InboundOrder{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestInboundOrderService_Delete_NotFound(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent order, got nil")
	}
}

func TestInboundOrderService_Confirm_Success(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return &model.InboundOrder{
				BaseModel:     model.BaseModel{ID: id},
				OrderNo:       "PO-2024-001",
				WarehouseID:   1,
				TotalQuantity: 100,
				Status:        0,
			}, nil
		},
		updateFunc: func(ctx context.Context, order *model.InboundOrder) error {
			return nil
		},
	}
	mockItemRepo := &mockInboundItemRepository{
		listByOrderIDFunc: func(ctx context.Context, orderID int64) ([]model.InboundItem, error) {
			return []model.InboundItem{
				{OrderID: orderID, ProductID: 1, Quantity: 50},
				{OrderID: orderID, ProductID: 2, Quantity: 50},
			}, nil
		},
	}
	mockInventorySvc := &mockInventoryServiceForInbound{
		adjustQuantityFunc: func(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
			return &model.Inventory{}, nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, mockItemRepo, mockInventorySvc, nil)

	order, err := svc.Confirm(context.Background(), 1)

	if err != nil {
		t.Fatalf("Confirm failed: %v", err)
	}
	if order.Status != 1 {
		t.Errorf("expected status 1 (completed), got %d", order.Status)
	}
}

func TestInboundOrderService_Confirm_AlreadyCompleted(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return &model.InboundOrder{
				BaseModel: model.BaseModel{ID: id},
				Status:    1,
			}, nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for already completed order, got nil")
	}
}

func TestInboundOrderService_Confirm_AlreadyCancelled(t *testing.T) {
	mockOrderRepo := &mockInboundOrderRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.InboundOrder, error) {
			return &model.InboundOrder{
				BaseModel: model.BaseModel{ID: id},
				Status:    2,
			}, nil
		},
	}

	svc := NewInboundOrderService(mockOrderRepo, nil, nil, nil)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for cancelled order, got nil")
	}
}

type mockInventoryServiceForInbound struct {
	adjustQuantityFunc func(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error)
}

func (m *mockInventoryServiceForInbound) AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
	if m.adjustQuantityFunc != nil {
		return m.adjustQuantityFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}
