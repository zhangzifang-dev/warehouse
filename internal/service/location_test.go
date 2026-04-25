package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockLocationRepository struct {
	createFunc              func(ctx context.Context, location *model.Location) error
	getByIDFunc             func(ctx context.Context, id int64) (*model.Location, error)
	getByWarehouseAndCodeFunc func(ctx context.Context, warehouseID int64, code string) (*model.Location, error)
	listFunc                func(ctx context.Context, page, pageSize int, warehouseID int64) ([]model.Location, int, error)
	updateFunc              func(ctx context.Context, location *model.Location) error
	deleteFunc              func(ctx context.Context, id int64) error
}

func (m *mockLocationRepository) Create(ctx context.Context, location *model.Location) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, location)
	}
	return errors.New("not implemented")
}

func (m *mockLocationRepository) GetByID(ctx context.Context, id int64) (*model.Location, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLocationRepository) GetByWarehouseAndCode(ctx context.Context, warehouseID int64, code string) (*model.Location, error) {
	if m.getByWarehouseAndCodeFunc != nil {
		return m.getByWarehouseAndCodeFunc(ctx, warehouseID, code)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLocationRepository) List(ctx context.Context, page, pageSize int, warehouseID int64) ([]model.Location, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, warehouseID)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockLocationRepository) Update(ctx context.Context, location *model.Location) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, location)
	}
	return errors.New("not implemented")
}

func (m *mockLocationRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

type mockWarehouseRepoForLocation struct {
	getByIDFunc func(ctx context.Context, id int64) (*model.Warehouse, error)
}

func (m *mockWarehouseRepoForLocation) Create(ctx context.Context, warehouse *model.Warehouse) error {
	return errors.New("not implemented")
}

func (m *mockWarehouseRepoForLocation) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseRepoForLocation) GetByCode(ctx context.Context, code string) (*model.Warehouse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockWarehouseRepoForLocation) List(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error) {
	return nil, 0, errors.New("not implemented")
}

func (m *mockWarehouseRepoForLocation) Update(ctx context.Context, warehouse *model.Warehouse) error {
	return errors.New("not implemented")
}

func (m *mockWarehouseRepoForLocation) Delete(ctx context.Context, id int64) error {
	return errors.New("not implemented")
}

func TestLocationService_Create_Success(t *testing.T) {
	createdLocation := &model.Location{}
	mockLocRepo := &mockLocationRepository{
		createFunc: func(ctx context.Context, location *model.Location) error {
			location.ID = 1
			createdLocation = location
			return nil
		},
		getByWarehouseAndCodeFunc: func(ctx context.Context, warehouseID int64, code string) (*model.Location, error) {
			return nil, errors.New("not found")
		},
	}
	mockWhRepo := &mockWarehouseRepoForLocation{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{BaseModel: model.BaseModel{ID: id}}, nil
		},
	}

	svc := NewLocationService(mockLocRepo, mockWhRepo)
	input := &CreateLocationInput{
		WarehouseID: 1,
		Zone:        "A",
		Shelf:       "01",
		Level:       "02",
		Position:    "03",
	}

	location, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if location == nil {
		t.Fatal("expected location, got nil")
	}
	if createdLocation.Code != "A-01-02-03" {
		t.Errorf("expected code 'A-01-02-03', got '%s'", createdLocation.Code)
	}
}

func TestLocationService_Create_WarehouseNotFound(t *testing.T) {
	mockLocRepo := &mockLocationRepository{}
	mockWhRepo := &mockWarehouseRepoForLocation{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewLocationService(mockLocRepo, mockWhRepo)
	input := &CreateLocationInput{
		WarehouseID: 999,
		Zone:        "A",
		Shelf:       "01",
		Level:       "02",
		Position:    "03",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for non-existent warehouse, got nil")
	}
}

func TestLocationService_Create_DuplicateCode(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		getByWarehouseAndCodeFunc: func(ctx context.Context, warehouseID int64, code string) (*model.Location, error) {
			return &model.Location{BaseModel: model.BaseModel{ID: 1}, Code: code}, nil
		},
	}
	mockWhRepo := &mockWarehouseRepoForLocation{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{BaseModel: model.BaseModel{ID: id}}, nil
		},
	}

	svc := NewLocationService(mockLocRepo, mockWhRepo)
	input := &CreateLocationInput{
		WarehouseID: 1,
		Zone:        "A",
		Shelf:       "01",
		Level:       "02",
		Position:    "03",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestLocationService_Create_DefaultStatus(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		createFunc: func(ctx context.Context, location *model.Location) error {
			return nil
		},
		getByWarehouseAndCodeFunc: func(ctx context.Context, warehouseID int64, code string) (*model.Location, error) {
			return nil, errors.New("not found")
		},
	}
	mockWhRepo := &mockWarehouseRepoForLocation{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Warehouse, error) {
			return &model.Warehouse{BaseModel: model.BaseModel{ID: id}}, nil
		},
	}

	svc := NewLocationService(mockLocRepo, mockWhRepo)
	input := &CreateLocationInput{
		WarehouseID: 1,
		Zone:        "A",
		Shelf:       "01",
		Level:       "02",
		Position:    "03",
	}

	location, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if location.Status != model.LocationStatusActive {
		t.Errorf("expected status %d, got %d", model.LocationStatusActive, location.Status)
	}
}

func TestLocationService_GetByID_Success(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return &model.Location{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				Code:        "A-01-02-03",
			}, nil
		},
	}

	svc := NewLocationService(mockLocRepo, nil)

	location, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if location == nil {
		t.Fatal("expected location, got nil")
	}
	if location.Code != "A-01-02-03" {
		t.Errorf("expected code 'A-01-02-03', got '%s'", location.Code)
	}
}

func TestLocationService_GetByID_NotFound(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewLocationService(mockLocRepo, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent location, got nil")
	}
}

func TestLocationService_List_Success(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		listFunc: func(ctx context.Context, page, pageSize int, warehouseID int64) ([]model.Location, int, error) {
			return []model.Location{
				{BaseModel: model.BaseModel{ID: 1}, WarehouseID: warehouseID, Code: "A-01-02-03"},
				{BaseModel: model.BaseModel{ID: 2}, WarehouseID: warehouseID, Code: "A-01-02-04"},
			}, 2, nil
		},
	}

	svc := NewLocationService(mockLocRepo, nil)

	result, err := svc.List(context.Background(), 1, 10, 1)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Locations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(result.Locations))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestLocationService_List_DefaultPagination(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		listFunc: func(ctx context.Context, page, pageSize int, warehouseID int64) ([]model.Location, int, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return []model.Location{}, 0, nil
		},
	}

	svc := NewLocationService(mockLocRepo, nil)

	_, err := svc.List(context.Background(), 0, 0, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestLocationService_Update_Success(t *testing.T) {
	updatedLocation := &model.Location{}
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return &model.Location{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				Zone:        "A",
				Shelf:       "01",
				Level:       "02",
				Position:    "03",
				Code:        "A-01-02-03",
			}, nil
		},
		updateFunc: func(ctx context.Context, location *model.Location) error {
			updatedLocation = location
			return nil
		},
	}

	svc := NewLocationService(mockLocRepo, nil)
	input := &UpdateLocationInput{
		Zone:     "B",
		Shelf:    "02",
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedLocation.Zone != "B" {
		t.Errorf("expected zone 'B', got '%s'", updatedLocation.Zone)
	}
	if updatedLocation.Code != "B-02-02-03" {
		t.Errorf("expected code 'B-02-02-03', got '%s'", updatedLocation.Code)
	}
}

func TestLocationService_Update_Status(t *testing.T) {
	newStatus := model.LocationStatusInactive
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return &model.Location{
				BaseModel:   model.BaseModel{ID: id},
				WarehouseID: 1,
				Code:        "A-01-02-03",
				Status:      model.LocationStatusActive,
			}, nil
		},
		updateFunc: func(ctx context.Context, location *model.Location) error {
			return nil
		},
	}

	svc := NewLocationService(mockLocRepo, nil)
	input := &UpdateLocationInput{
		Status: &newStatus,
	}

	location, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if location.Status != model.LocationStatusInactive {
		t.Errorf("expected status %d, got %d", model.LocationStatusInactive, location.Status)
	}
}

func TestLocationService_Update_NotFound(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewLocationService(mockLocRepo, nil)
	input := &UpdateLocationInput{Zone: "B"}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent location, got nil")
	}
}

func TestLocationService_Delete_Success(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return &model.Location{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewLocationService(mockLocRepo, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestLocationService_Delete_NotFound(t *testing.T) {
	mockLocRepo := &mockLocationRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Location, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewLocationService(mockLocRepo, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent location, got nil")
	}
}
