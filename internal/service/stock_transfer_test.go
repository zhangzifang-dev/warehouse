package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockStockTransferRepository struct {
	createFunc       func(ctx context.Context, transfer *model.StockTransfer) error
	getByIDFunc      func(ctx context.Context, id int64) (*model.StockTransfer, error)
	getByOrderNoFunc func(ctx context.Context, orderNo string) (*model.StockTransfer, error)
	listFunc         func(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) ([]model.StockTransfer, int, error)
	updateFunc       func(ctx context.Context, transfer *model.StockTransfer) error
	deleteFunc       func(ctx context.Context, id int64) error
}

func (m *mockStockTransferRepository) Create(ctx context.Context, transfer *model.StockTransfer) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, transfer)
	}
	return errors.New("not implemented")
}

func (m *mockStockTransferRepository) GetByID(ctx context.Context, id int64) (*model.StockTransfer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.StockTransfer, error) {
	if m.getByOrderNoFunc != nil {
		return m.getByOrderNoFunc(ctx, orderNo)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferRepository) List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) ([]model.StockTransfer, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, fromWarehouseID, toWarehouseID, status)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockStockTransferRepository) Update(ctx context.Context, transfer *model.StockTransfer) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, transfer)
	}
	return errors.New("not implemented")
}

func (m *mockStockTransferRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

type mockStockTransferItemRepository struct {
	createFunc          func(ctx context.Context, item *model.StockTransferItem) error
	listByTransferIDFunc func(ctx context.Context, transferID int64) ([]model.StockTransferItem, error)
	updateFunc          func(ctx context.Context, item *model.StockTransferItem) error
	deleteFunc          func(ctx context.Context, id int64) error
}

func (m *mockStockTransferItemRepository) Create(ctx context.Context, item *model.StockTransferItem) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, item)
	}
	return errors.New("not implemented")
}

func (m *mockStockTransferItemRepository) ListByTransferID(ctx context.Context, transferID int64) ([]model.StockTransferItem, error) {
	if m.listByTransferIDFunc != nil {
		return m.listByTransferIDFunc(ctx, transferID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockStockTransferItemRepository) Update(ctx context.Context, item *model.StockTransferItem) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, item)
	}
	return errors.New("not implemented")
}

func (m *mockStockTransferItemRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

type mockInventoryServiceForTransfer struct {
	checkStockFunc              func(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error)
	adjustQuantityFunc          func(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error)
	getByWarehouseAndProductFunc func(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error)
}

func (m *mockInventoryServiceForTransfer) CheckStock(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
	if m.checkStockFunc != nil {
		return m.checkStockFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryServiceForTransfer) AdjustQuantity(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
	if m.adjustQuantityFunc != nil {
		return m.adjustQuantityFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryServiceForTransfer) GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
	if m.getByWarehouseAndProductFunc != nil {
		return m.getByWarehouseAndProductFunc(ctx, warehouseID, productID, batchNo)
	}
	return nil, errors.New("not implemented")
}

func TestStockTransferService_Create_Success(t *testing.T) {
	createdTransfer := &model.StockTransfer{}
	mockTransferRepo := &mockStockTransferRepository{
		createFunc: func(ctx context.Context, transfer *model.StockTransfer) error {
			transfer.ID = 1
			createdTransfer = transfer
			return nil
		},
	}
	mockItemRepo := &mockStockTransferItemRepository{}

	svc := NewStockTransferService(mockTransferRepo, mockItemRepo, nil)
	input := &CreateStockTransferInput{
		OrderNo:         "ST-2024-001",
		FromWarehouseID: 1,
		ToWarehouseID:   2,
		TotalQuantity:   100,
	}

	transfer, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if transfer == nil {
		t.Fatal("expected transfer, got nil")
	}
	if createdTransfer.OrderNo != "ST-2024-001" {
		t.Errorf("expected OrderNo 'ST-2024-001', got '%s'", createdTransfer.OrderNo)
	}
}

func TestStockTransferService_Create_MissingFromWarehouseID(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{}
	mockItemRepo := &mockStockTransferItemRepository{}

	svc := NewStockTransferService(mockTransferRepo, mockItemRepo, nil)
	input := &CreateStockTransferInput{
		OrderNo:         "ST-2024-001",
		ToWarehouseID:   2,
		TotalQuantity:   100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing from warehouse ID, got nil")
	}
}

func TestStockTransferService_Create_MissingToWarehouseID(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{}
	mockItemRepo := &mockStockTransferItemRepository{}

	svc := NewStockTransferService(mockTransferRepo, mockItemRepo, nil)
	input := &CreateStockTransferInput{
		OrderNo:         "ST-2024-001",
		FromWarehouseID: 1,
		TotalQuantity:   100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing to warehouse ID, got nil")
	}
}

func TestStockTransferService_Create_SameWarehouse(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{}
	mockItemRepo := &mockStockTransferItemRepository{}

	svc := NewStockTransferService(mockTransferRepo, mockItemRepo, nil)
	input := &CreateStockTransferInput{
		OrderNo:         "ST-2024-001",
		FromWarehouseID: 1,
		ToWarehouseID:   1,
		TotalQuantity:   100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for same warehouse, got nil")
	}
}

func TestStockTransferService_GetByID_Success(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel:       model.BaseModel{ID: id},
				OrderNo:         "ST-2024-001",
				FromWarehouseID: 1,
				ToWarehouseID:   2,
				TotalQuantity:   100,
			}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	transfer, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if transfer == nil {
		t.Fatal("expected transfer, got nil")
	}
	if transfer.TotalQuantity != 100 {
		t.Errorf("expected total quantity 100, got %f", transfer.TotalQuantity)
	}
}

func TestStockTransferService_GetByID_NotFound(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent transfer, got nil")
	}
}

func TestStockTransferService_GetByOrderNo_Success(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByOrderNoFunc: func(ctx context.Context, orderNo string) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel:       model.BaseModel{ID: 1},
				OrderNo:         orderNo,
				FromWarehouseID: 1,
				ToWarehouseID:   2,
			}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	transfer, err := svc.GetByOrderNo(context.Background(), "ST-2024-001")

	if err != nil {
		t.Fatalf("GetByOrderNo failed: %v", err)
	}
	if transfer == nil {
		t.Fatal("expected transfer, got nil")
	}
}

func TestStockTransferService_List_Success(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		listFunc: func(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) ([]model.StockTransfer, int, error) {
			return []model.StockTransfer{
				{BaseModel: model.BaseModel{ID: 1}, OrderNo: "ST-2024-001", FromWarehouseID: 1, ToWarehouseID: 2, TotalQuantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, OrderNo: "ST-2024-002", FromWarehouseID: 1, ToWarehouseID: 3, TotalQuantity: 200},
			}, 2, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	result, err := svc.List(context.Background(), 1, 10, 0, 0, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Transfers) != 2 {
		t.Errorf("expected 2 transfers, got %d", len(result.Transfers))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestStockTransferService_Update_Success(t *testing.T) {
	updatedTransfer := &model.StockTransfer{}
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel:       model.BaseModel{ID: id},
				OrderNo:         "ST-2024-001",
				FromWarehouseID: 1,
				ToWarehouseID:   2,
				TotalQuantity:   100,
				Status:          0,
			}, nil
		},
		updateFunc: func(ctx context.Context, transfer *model.StockTransfer) error {
			updatedTransfer = transfer
			return nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)
	input := &UpdateStockTransferInput{
		TotalQuantity: floatPtrST(200),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedTransfer.TotalQuantity != 200 {
		t.Errorf("expected total quantity 200, got %f", updatedTransfer.TotalQuantity)
	}
}

func TestStockTransferService_Update_SameWarehouse(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel:       model.BaseModel{ID: id},
				FromWarehouseID: 1,
				ToWarehouseID:   2,
			}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)
	input := &UpdateStockTransferInput{
		ToWarehouseID: int64PtrST(1),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for same warehouse, got nil")
	}
}

func TestStockTransferService_Delete_Success(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestStockTransferService_Delete_NotFound(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent transfer, got nil")
	}
}

func TestStockTransferService_Confirm_Success(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel:       model.BaseModel{ID: id},
				OrderNo:         "ST-2024-001",
				FromWarehouseID: 1,
				ToWarehouseID:   2,
				TotalQuantity:   100,
				Status:          0,
			}, nil
		},
		updateFunc: func(ctx context.Context, transfer *model.StockTransfer) error {
			return nil
		},
	}
	mockItemRepo := &mockStockTransferItemRepository{
		listByTransferIDFunc: func(ctx context.Context, transferID int64) ([]model.StockTransferItem, error) {
			return []model.StockTransferItem{
				{TransferID: transferID, ProductID: 1, Quantity: 50},
				{TransferID: transferID, ProductID: 2, Quantity: 50},
			}, nil
		},
	}
	mockInventorySvc := &mockInventoryServiceForTransfer{
		checkStockFunc: func(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
			return &CheckStockResult{Available: true}, nil
		},
		getByWarehouseAndProductFunc: func(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
			return &model.Inventory{BaseModel: model.BaseModel{ID: 1}, Quantity: 100}, nil
		},
		adjustQuantityFunc: func(ctx context.Context, input *AdjustQuantityInput) (*model.Inventory, error) {
			return &model.Inventory{}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, mockItemRepo, mockInventorySvc)

	transfer, err := svc.Confirm(context.Background(), 1)

	if err != nil {
		t.Fatalf("Confirm failed: %v", err)
	}
	if transfer.Status != 1 {
		t.Errorf("expected status 1 (completed), got %d", transfer.Status)
	}
}

func TestStockTransferService_Confirm_AlreadyCompleted(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel: model.BaseModel{ID: id},
				Status:    1,
			}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for already completed transfer, got nil")
	}
}

func TestStockTransferService_Confirm_AlreadyCancelled(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel: model.BaseModel{ID: id},
				Status:    2,
			}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, nil, nil)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for cancelled transfer, got nil")
	}
}

func TestStockTransferService_Confirm_InsufficientStock(t *testing.T) {
	mockTransferRepo := &mockStockTransferRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.StockTransfer, error) {
			return &model.StockTransfer{
				BaseModel:       model.BaseModel{ID: id},
				FromWarehouseID: 1,
				ToWarehouseID:   2,
				Status:          0,
			}, nil
		},
	}
	mockItemRepo := &mockStockTransferItemRepository{
		listByTransferIDFunc: func(ctx context.Context, transferID int64) ([]model.StockTransferItem, error) {
			return []model.StockTransferItem{
				{TransferID: transferID, ProductID: 1, Quantity: 50},
			}, nil
		},
	}
	mockInventorySvc := &mockInventoryServiceForTransfer{
		checkStockFunc: func(ctx context.Context, input *CheckStockInput) (*CheckStockResult, error) {
			return &CheckStockResult{Available: false}, nil
		},
	}

	svc := NewStockTransferService(mockTransferRepo, mockItemRepo, mockInventorySvc)

	_, err := svc.Confirm(context.Background(), 1)

	if err == nil {
		t.Error("expected error for insufficient stock, got nil")
	}
}

func floatPtrST(f float64) *float64 {
	return &f
}

func int64PtrST(i int64) *int64 {
	return &i
}
