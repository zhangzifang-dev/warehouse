package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockWarehouseRepository struct {
	createFunc    func(ctx context.Context, warehouse *model.Warehouse) error
	getByIDFunc   func(ctx context.Context, id int64) (*model.Warehouse, error)
	getByCodeFunc func(ctx context.Context, code string) (*model.Warehouse, error)
	listFunc      func(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error)
	updateFunc    func(ctx context.Context, warehouse *model.Warehouse) error
	deleteFunc    func(ctx context.Context, id int64) error
}

func (m *mockWarehouseRepository) Create(ctx context.Context, warehouse *model.Warehouse) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, warehouse)
	}
	return errors.New("not implemented")
}

func (m *mockWarehouseRepository) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseRepository) GetByCode(ctx context.Context, code string) (*model.Warehouse, error) {
	if m.getByCodeFunc != nil {
		return m.getByCodeFunc(ctx, code)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseRepository) List(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockWarehouseRepository) Update(ctx context.Context, warehouse *model.Warehouse) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, warehouse)
	}
	return errors.New("not implemented")
}

func (m *mockWarehouseRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestWarehouseService_Create_Success(t *testing.T) {
	createdWarehouse := &model.Warehouse{}
	mockRepo := &mockWarehouseRepository{
		createFunc: func(ctx context.Context, warehouse *model.Warehouse) error {
			warehouse.ID = 1
			createdWarehouse = warehouse
			return nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Warehouse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewWarehouseService(mockRepo)
	input := &CreateWarehouseInput{
		Name:    "Main Warehouse",
		Code:    "WH001",
		Address: "123 Main St",
		Contact: "John Doe",
		Phone:   "1234567890",
	}

	warehouse, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if warehouse == nil {
		t.Fatal("expected warehouse, got nil")
	}
	if createdWarehouse.Name != "Main Warehouse" {
		t.Errorf("expected name 'Main Warehouse', got '%s'", createdWarehouse.Name)
	}
}

func TestWarehouseService_Create_DefaultStatus(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		createFunc: func(ctx context.Context, warehouse *model.Warehouse) error {
			return nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Warehouse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewWarehouseService(mockRepo)
	input := &CreateWarehouseInput{
		Name: "Test Warehouse",
		Code: "WH002",
	}

	warehouse, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if warehouse.Status != model.WarehouseStatusActive {
		t.Errorf("expected status %d, got %d", model.WarehouseStatusActive, warehouse.Status)
	}
}

func TestWarehouseService_Create_DuplicateCode(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		getByCodeFunc: func(ctx context.Context, code string) (*model.Warehouse, error) {
			return &model.Warehouse{BaseModel: model.BaseModel{ID: 1}, Code: code}, nil
		},
	}

	svc := NewWarehouseService(mockRepo)
	input := &CreateWarehouseInput{
		Name: "Test Warehouse",
		Code: "WH001",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestWarehouseService_GetByID_Success(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Main Warehouse",
				Code:      "WH001",
			}, nil
		},
	}

	svc := NewWarehouseService(mockRepo)

	warehouse, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if warehouse == nil {
		t.Fatal("expected warehouse, got nil")
	}
	if warehouse.Name != "Main Warehouse" {
		t.Errorf("expected name 'Main Warehouse', got '%s'", warehouse.Name)
	}
}

func TestWarehouseService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewWarehouseService(mockRepo)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent warehouse, got nil")
	}
}

func TestWarehouseService_List_Success(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error) {
			return []model.Warehouse{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Warehouse 1", Code: "WH001"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Warehouse 2", Code: "WH002"},
			}, 2, nil
		},
	}

	svc := NewWarehouseService(mockRepo)

	result, err := svc.List(context.Background(), 1, 10)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Warehouses) != 2 {
		t.Errorf("expected 2 warehouses, got %d", len(result.Warehouses))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestWarehouseService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return []model.Warehouse{}, 0, nil
		},
	}

	svc := NewWarehouseService(mockRepo)

	_, err := svc.List(context.Background(), 0, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestWarehouseService_List_MaxPageSize(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error) {
			if pageSize > 100 {
				t.Errorf("expected pageSize <= 100, got %d", pageSize)
			}
			return []model.Warehouse{}, 0, nil
		},
	}

	svc := NewWarehouseService(mockRepo)

	_, err := svc.List(context.Background(), 1, 200)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestWarehouseService_Update_Success(t *testing.T) {
	updatedWarehouse := &model.Warehouse{}
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Old Name",
				Code:      "WH001",
			}, nil
		},
		updateFunc: func(ctx context.Context, warehouse *model.Warehouse) error {
			updatedWarehouse = warehouse
			return nil
		},
	}

	svc := NewWarehouseService(mockRepo)
	input := &UpdateWarehouseInput{
		Name:    "New Name",
		Address: "New Address",
	}

	warehouse, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if warehouse.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", updatedWarehouse.Name)
	}
}

func TestWarehouseService_Update_Status(t *testing.T) {
	newStatus := model.WarehouseStatusDisabled
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Warehouse",
				Code:      "WH001",
				Status:    model.WarehouseStatusActive,
			}, nil
		},
		updateFunc: func(ctx context.Context, warehouse *model.Warehouse) error {
			return nil
		},
	}

	svc := NewWarehouseService(mockRepo)
	input := &UpdateWarehouseInput{
		Status: &newStatus,
	}

	warehouse, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if warehouse.Status != model.WarehouseStatusDisabled {
		t.Errorf("expected status %d, got %d", model.WarehouseStatusDisabled, warehouse.Status)
	}
}

func TestWarehouseService_Update_NotFound(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewWarehouseService(mockRepo)
	input := &UpdateWarehouseInput{Name: "New Name"}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent warehouse, got nil")
	}
}

func TestWarehouseService_Delete_Success(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewWarehouseService(mockRepo)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestWarehouseService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockWarehouseRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewWarehouseService(mockRepo)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent warehouse, got nil")
	}
}
