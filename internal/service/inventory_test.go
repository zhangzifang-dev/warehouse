package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type mockInventoryRepository struct {
	createFunc                 func(ctx context.Context, inventory *model.Inventory) error
	getByIDFunc                func(ctx context.Context, id int64) (*model.Inventory, error)
	listFunc                   func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error)
	updateFunc                 func(ctx context.Context, inventory *model.Inventory) error
	deleteFunc                 func(ctx context.Context, id int64) error
	getByWarehouseAndProductFunc func(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error)
	updateQuantityFunc         func(ctx context.Context, id int64, quantity float64) error
}

func (m *mockInventoryRepository) Create(ctx context.Context, inventory *model.Inventory) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, inventory)
	}
	return errors.New("not implemented")
}

func (m *mockInventoryRepository) GetByID(ctx context.Context, id int64) (*model.Inventory, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryRepository) List(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockInventoryRepository) Update(ctx context.Context, inventory *model.Inventory) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, inventory)
	}
	return errors.New("not implemented")
}

func (m *mockInventoryRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockInventoryRepository) GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
	if m.getByWarehouseAndProductFunc != nil {
		return m.getByWarehouseAndProductFunc(ctx, warehouseID, productID, batchNo)
	}
	return nil, errors.New("not implemented")
}

func (m *mockInventoryRepository) UpdateQuantity(ctx context.Context, id int64, quantity float64) error {
	if m.updateQuantityFunc != nil {
		return m.updateQuantityFunc(ctx, id, quantity)
	}
	return errors.New("not implemented")
}

func TestInventoryService_Create_Success(t *testing.T) {
	createdInventory := &model.Inventory{}
	mockRepo := &mockInventoryRepository{
		createFunc: func(ctx context.Context, inventory *model.Inventory) error {
			inventory.ID = 1
			createdInventory = inventory
			return nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &CreateInventoryInput{
		WarehouseID: 1,
		ProductID:   1,
		Quantity:    100,
		BatchNo:     "BATCH001",
	}

	inventory, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if inventory == nil {
		t.Fatal("expected inventory, got nil")
	}
	if createdInventory.WarehouseID != 1 {
		t.Errorf("expected WarehouseID 1, got %d", createdInventory.WarehouseID)
	}
	if createdInventory.ProductID != 1 {
		t.Errorf("expected ProductID 1, got %d", createdInventory.ProductID)
	}
}

func TestInventoryService_Create_MissingWarehouseID(t *testing.T) {
	mockRepo := &mockInventoryRepository{}

	svc := NewInventoryService(mockRepo, nil)
	input := &CreateInventoryInput{
		ProductID: 1,
		Quantity:  100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing warehouse ID, got nil")
	}
}

func TestInventoryService_Create_MissingProductID(t *testing.T) {
	mockRepo := &mockInventoryRepository{}

	svc := NewInventoryService(mockRepo, nil)
	input := &CreateInventoryInput{
		WarehouseID: 1,
		Quantity:    100,
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing product ID, got nil")
	}
}

func TestInventoryService_GetByID_Success(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    100,
			}, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	inventory, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if inventory == nil {
		t.Fatal("expected inventory, got nil")
	}
	if inventory.Quantity != 100 {
		t.Errorf("expected quantity 100, got %f", inventory.Quantity)
	}
}

func TestInventoryService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent inventory, got nil")
	}
}

func TestInventoryService_List_Success(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, ProductID: 1, Quantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, WarehouseID: 1, ProductID: 2, Quantity: 200},
			}, 2, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	result, err := svc.List(context.Background(), &model.InventoryQueryFilter{Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Inventories) != 2 {
		t.Errorf("expected 2 inventories, got %d", len(result.Inventories))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestInventoryService_List_WithProductName(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.ProductName != "测试商品" {
				t.Errorf("expected ProductName '测试商品', got '%s'", filter.ProductName)
			}
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, ProductID: 1, Quantity: 100},
			}, 1, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{
		Page:        1,
		PageSize:    10,
		ProductName: "测试商品",
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_List_WithQuantityMin(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.QuantityMin == nil || *filter.QuantityMin != 50 {
				t.Errorf("expected QuantityMin 50, got %v", filter.QuantityMin)
			}
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, Quantity: 100},
			}, 1, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	minQty := 50.0

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{
		Page:        1,
		PageSize:    10,
		QuantityMin: &minQty,
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_List_WithQuantityMax(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.QuantityMax == nil || *filter.QuantityMax != 200 {
				t.Errorf("expected QuantityMax 200, got %v", filter.QuantityMax)
			}
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, Quantity: 100},
			}, 1, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	maxQty := 200.0

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{
		Page:        1,
		PageSize:    10,
		QuantityMax: &maxQty,
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_List_WithQuantityRange(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.QuantityMin == nil || *filter.QuantityMin != 50 {
				t.Errorf("expected QuantityMin 50, got %v", filter.QuantityMin)
			}
			if filter.QuantityMax == nil || *filter.QuantityMax != 200 {
				t.Errorf("expected QuantityMax 200, got %v", filter.QuantityMax)
			}
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, Quantity: 100},
				{BaseModel: model.BaseModel{ID: 2}, Quantity: 150},
			}, 2, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	minQty := 50.0
	maxQty := 200.0

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{
		Page:        1,
		PageSize:    10,
		QuantityMin: &minQty,
		QuantityMax: &maxQty,
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_List_WithBatchNo(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.BatchNo != "BATCH001" {
				t.Errorf("expected BatchNo 'BATCH001', got '%s'", filter.BatchNo)
			}
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, BatchNo: "BATCH001", Quantity: 100},
			}, 1, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{
		Page:     1,
		PageSize: 10,
		BatchNo:  "BATCH001",
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_List_WithAllFilters(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.ProductName != "测试商品" {
				t.Errorf("expected ProductName '测试商品', got '%s'", filter.ProductName)
			}
			if filter.QuantityMin == nil || *filter.QuantityMin != 50 {
				t.Errorf("expected QuantityMin 50, got %v", filter.QuantityMin)
			}
			if filter.QuantityMax == nil || *filter.QuantityMax != 200 {
				t.Errorf("expected QuantityMax 200, got %v", filter.QuantityMax)
			}
			if filter.BatchNo != "BATCH001" {
				t.Errorf("expected BatchNo 'BATCH001', got '%s'", filter.BatchNo)
			}
			return []model.Inventory{
				{BaseModel: model.BaseModel{ID: 1}, ProductID: 1, BatchNo: "BATCH001", Quantity: 100},
			}, 1, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	minQty := 50.0
	maxQty := 200.0

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{
		Page:        1,
		PageSize:    10,
		ProductName: "测试商品",
		QuantityMin: &minQty,
		QuantityMax: &maxQty,
		BatchNo:     "BATCH001",
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		listFunc: func(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
			if filter.Page != 1 {
				t.Errorf("expected page 1, got %d", filter.Page)
			}
			if filter.PageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", filter.PageSize)
			}
			return []model.Inventory{}, 0, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	_, err := svc.List(context.Background(), &model.InventoryQueryFilter{Page: 0, PageSize: 0})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestInventoryService_Update_Success(t *testing.T) {
	updatedInventory := &model.Inventory{}
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    100,
			}, nil
		},
		updateFunc: func(ctx context.Context, inventory *model.Inventory) error {
			updatedInventory = inventory
			return nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &UpdateInventoryInput{
		Quantity: invFloatPtr(200),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedInventory.Quantity != 200 {
		t.Errorf("expected quantity 200, got %f", updatedInventory.Quantity)
	}
}

func TestInventoryService_Update_NotFound(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &UpdateInventoryInput{Quantity: invFloatPtr(200)}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent inventory, got nil")
	}
}

func TestInventoryService_Delete_Success(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return &model.Inventory{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestInventoryService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInventoryService(mockRepo, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent inventory, got nil")
	}
}

func TestInventoryService_AdjustQuantity_Increase(t *testing.T) {
	updatedInventory := &model.Inventory{}
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    100,
			}, nil
		},
		updateFunc: func(ctx context.Context, inventory *model.Inventory) error {
			updatedInventory = inventory
			return nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &AdjustQuantityInput{
		InventoryID: 1,
		Quantity:    50,
	}

	inventory, err := svc.AdjustQuantity(context.Background(), input)

	if err != nil {
		t.Fatalf("AdjustQuantity failed: %v", err)
	}
	if inventory.Quantity != 150 {
		t.Errorf("expected quantity 150, got %f", inventory.Quantity)
	}
	if updatedInventory.Quantity != 150 {
		t.Errorf("expected updated quantity 150, got %f", updatedInventory.Quantity)
	}
}

func TestInventoryService_AdjustQuantity_Decrease(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    100,
			}, nil
		},
		updateFunc: func(ctx context.Context, inventory *model.Inventory) error {
			return nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &AdjustQuantityInput{
		InventoryID: 1,
		Quantity:    -50,
	}

	inventory, err := svc.AdjustQuantity(context.Background(), input)

	if err != nil {
		t.Fatalf("AdjustQuantity failed: %v", err)
	}
	if inventory.Quantity != 50 {
		t.Errorf("expected quantity 50, got %f", inventory.Quantity)
	}
}

func TestInventoryService_AdjustQuantity_InsufficientStock(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				ProductID:   1,
				Quantity:    50,
			}, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &AdjustQuantityInput{
		InventoryID: 1,
		Quantity:    -100,
	}

	_, err := svc.AdjustQuantity(context.Background(), input)

	if err == nil {
		t.Error("expected error for insufficient stock, got nil")
	}
	appErr, ok := apperrors.GetAppError(err)
	if !ok || appErr.Code != apperrors.CodeInsufficientStock {
		t.Errorf("expected CodeInsufficientStock error, got %v", err)
	}
}

func TestInventoryService_AdjustQuantity_NotFound(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Inventory, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &AdjustQuantityInput{
		InventoryID: 999,
		Quantity:    50,
	}

	_, err := svc.AdjustQuantity(context.Background(), input)

	if err == nil {
		t.Error("expected error for non-existent inventory, got nil")
	}
}

func TestInventoryService_AdjustQuantity_InvalidID(t *testing.T) {
	mockRepo := &mockInventoryRepository{}

	svc := NewInventoryService(mockRepo, nil)
	input := &AdjustQuantityInput{
		InventoryID: 0,
		Quantity:    50,
	}

	_, err := svc.AdjustQuantity(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid inventory ID, got nil")
	}
}

func TestInventoryService_CheckStock_Available(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByWarehouseAndProductFunc: func(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: 1},
				WarehouseID: warehouseID,
				ProductID:   productID,
				Quantity:    100,
				BatchNo:     batchNo,
			}, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &CheckStockInput{
		WarehouseID: 1,
		ProductID:   1,
		BatchNo:     "BATCH001",
		Quantity:    50,
	}

	result, err := svc.CheckStock(context.Background(), input)

	if err != nil {
		t.Fatalf("CheckStock failed: %v", err)
	}
	if !result.Available {
		t.Error("expected stock to be available")
	}
	if result.CurrentStock != 100 {
		t.Errorf("expected current stock 100, got %f", result.CurrentStock)
	}
}

func TestInventoryService_CheckStock_NotAvailable(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByWarehouseAndProductFunc: func(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
			return &model.Inventory{
				BaseModel:   model.BaseModel{ID: 1},
				WarehouseID: warehouseID,
				ProductID:   productID,
				Quantity:    30,
				BatchNo:     batchNo,
			}, nil
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &CheckStockInput{
		WarehouseID: 1,
		ProductID:   1,
		BatchNo:     "BATCH001",
		Quantity:    50,
	}

	result, err := svc.CheckStock(context.Background(), input)

	if err != nil {
		t.Fatalf("CheckStock failed: %v", err)
	}
	if result.Available {
		t.Error("expected stock to be unavailable")
	}
}

func TestInventoryService_CheckStock_NotFound(t *testing.T) {
	mockRepo := &mockInventoryRepository{
		getByWarehouseAndProductFunc: func(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewInventoryService(mockRepo, nil)
	input := &CheckStockInput{
		WarehouseID: 1,
		ProductID:   1,
		BatchNo:     "BATCH001",
		Quantity:    50,
	}

	result, err := svc.CheckStock(context.Background(), input)

	if err != nil {
		t.Fatalf("CheckStock failed: %v", err)
	}
	if result.Available {
		t.Error("expected stock to be unavailable when inventory not found")
	}
	if result.CurrentStock != 0 {
		t.Errorf("expected current stock 0, got %f", result.CurrentStock)
	}
}

func TestInventoryService_CheckStock_MissingWarehouseID(t *testing.T) {
	mockRepo := &mockInventoryRepository{}

	svc := NewInventoryService(mockRepo, nil)
	input := &CheckStockInput{
		ProductID: 1,
		Quantity:  50,
	}

	_, err := svc.CheckStock(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing warehouse ID, got nil")
	}
}

func TestInventoryService_CheckStock_MissingProductID(t *testing.T) {
	mockRepo := &mockInventoryRepository{}

	svc := NewInventoryService(mockRepo, nil)
	input := &CheckStockInput{
		WarehouseID: 1,
		Quantity:    50,
	}

	_, err := svc.CheckStock(context.Background(), input)

	if err == nil {
		t.Error("expected error for missing product ID, got nil")
	}
}

func TestInventoryService_CheckStock_InvalidQuantity(t *testing.T) {
	mockRepo := &mockInventoryRepository{}

	svc := NewInventoryService(mockRepo, nil)
	input := &CheckStockInput{
		WarehouseID: 1,
		ProductID:   1,
		Quantity:    0,
	}

	_, err := svc.CheckStock(context.Background(), input)

	if err == nil {
		t.Error("expected error for invalid quantity, got nil")
	}
}

func invFloatPtr(f float64) *float64 {
	return &f
}

func invIntPtr(i int64) *int64 {
	return &i
}

func invStrPtr(s string) *string {
	return &s
}
