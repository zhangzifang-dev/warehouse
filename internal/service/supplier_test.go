package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
	"warehouse/internal/repository"
)

type mockSupplierRepository struct {
	createFunc    func(ctx context.Context, supplier *model.Supplier) error
	getByIDFunc   func(ctx context.Context, id int64) (*model.Supplier, error)
	getByCodeFunc func(ctx context.Context, code string) (*model.Supplier, error)
	listFunc      func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error)
	updateFunc    func(ctx context.Context, supplier *model.Supplier) error
	deleteFunc    func(ctx context.Context, id int64) error
}

func (m *mockSupplierRepository) Create(ctx context.Context, supplier *model.Supplier) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, supplier)
	}
	return errors.New("not implemented")
}

func (m *mockSupplierRepository) GetByID(ctx context.Context, id int64) (*model.Supplier, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSupplierRepository) GetByCode(ctx context.Context, code string) (*model.Supplier, error) {
	if m.getByCodeFunc != nil {
		return m.getByCodeFunc(ctx, code)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSupplierRepository) List(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockSupplierRepository) Update(ctx context.Context, supplier *model.Supplier) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, supplier)
	}
	return errors.New("not implemented")
}

func (m *mockSupplierRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestSupplierService_Create_Success(t *testing.T) {
	createdSupplier := &model.Supplier{}
	mockRepo := &mockSupplierRepository{
		createFunc: func(ctx context.Context, supplier *model.Supplier) error {
			supplier.ID = 1
			createdSupplier = supplier
			return nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &CreateSupplierInput{
		Name:    "Test Supplier",
		Code:    "SUP001",
		Contact: "John Doe",
		Phone:   "1234567890",
	}

	supplier, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if supplier == nil {
		t.Fatal("expected supplier, got nil")
	}
	if createdSupplier.Name != "Test Supplier" {
		t.Errorf("expected name 'Test Supplier', got '%s'", createdSupplier.Name)
	}
}

func TestSupplierService_Create_EmptyName(t *testing.T) {
	mockRepo := &mockSupplierRepository{}

	svc := NewSupplierService(mockRepo, nil)
	input := &CreateSupplierInput{
		Name: "",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty name, got nil")
	}
}

func TestSupplierService_Create_DuplicateCode(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByCodeFunc: func(ctx context.Context, code string) (*model.Supplier, error) {
			return &model.Supplier{BaseModel: model.BaseModel{ID: 1}, Code: code}, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &CreateSupplierInput{
		Name: "Test Supplier",
		Code: "SUP001",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestSupplierService_Create_DefaultStatus(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		createFunc: func(ctx context.Context, supplier *model.Supplier) error {
			return nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &CreateSupplierInput{
		Name: "Test Supplier",
	}

	supplier, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if supplier.Status != model.SupplierStatusActive {
		t.Errorf("expected status %d, got %d", model.SupplierStatusActive, supplier.Status)
	}
}

func TestSupplierService_GetByID_Success(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return &model.Supplier{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Supplier",
			}, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	supplier, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if supplier == nil {
		t.Fatal("expected supplier, got nil")
	}
	if supplier.Name != "Test Supplier" {
		t.Errorf("expected name 'Test Supplier', got '%s'", supplier.Name)
	}
}

func TestSupplierService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent supplier, got nil")
	}
}

func TestSupplierService_List_Success(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			return []model.Supplier{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Supplier A"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Supplier B"},
			}, 2, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	result, err := svc.List(context.Background(), &SupplierQueryFilter{Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Suppliers) != 2 {
		t.Errorf("expected 2 suppliers, got %d", len(result.Suppliers))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestSupplierService_List_WithKeyword(t *testing.T) {
	receivedFilter := &repository.SupplierQueryFilter{}
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			receivedFilter = filter
			return []model.Supplier{{BaseModel: model.BaseModel{ID: 1}, Name: "Test Supplier"}}, 1, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	_, err := svc.List(context.Background(), &SupplierQueryFilter{Name: "test", Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedFilter.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", receivedFilter.Name)
	}
}

func TestSupplierService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			if filter.Page != 1 {
				t.Errorf("expected page 1, got %d", filter.Page)
			}
			if filter.PageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", filter.PageSize)
			}
			return []model.Supplier{}, 0, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	_, err := svc.List(context.Background(), &SupplierQueryFilter{Page: 0, PageSize: 0})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestSupplierService_List_WithCodeFilter(t *testing.T) {
	receivedFilter := &repository.SupplierQueryFilter{}
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			receivedFilter = filter
			return []model.Supplier{{BaseModel: model.BaseModel{ID: 1}, Code: "SUP001"}}, 1, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	_, err := svc.List(context.Background(), &SupplierQueryFilter{Code: "SUP", Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedFilter.Code != "SUP" {
		t.Errorf("expected code 'SUP', got '%s'", receivedFilter.Code)
	}
}

func TestSupplierService_List_WithContactFilter(t *testing.T) {
	receivedFilter := &repository.SupplierQueryFilter{}
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			receivedFilter = filter
			return []model.Supplier{{BaseModel: model.BaseModel{ID: 1}, Contact: "John"}}, 1, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	_, err := svc.List(context.Background(), &SupplierQueryFilter{Contact: "John", Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedFilter.Contact != "John" {
		t.Errorf("expected contact 'John', got '%s'", receivedFilter.Contact)
	}
}

func TestSupplierService_List_WithPhoneFilter(t *testing.T) {
	receivedFilter := &repository.SupplierQueryFilter{}
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			receivedFilter = filter
			return []model.Supplier{{BaseModel: model.BaseModel{ID: 1}, Phone: "123456"}}, 1, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	_, err := svc.List(context.Background(), &SupplierQueryFilter{Phone: "123", Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedFilter.Phone != "123" {
		t.Errorf("expected phone '123', got '%s'", receivedFilter.Phone)
	}
}

func TestSupplierService_List_WithStatusFilter(t *testing.T) {
	receivedFilter := &repository.SupplierQueryFilter{}
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			receivedFilter = filter
			return []model.Supplier{{BaseModel: model.BaseModel{ID: 1}, Status: 1}}, 1, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	status := 1

	_, err := svc.List(context.Background(), &SupplierQueryFilter{Status: &status, Page: 1, PageSize: 10})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedFilter.Status == nil || *receivedFilter.Status != 1 {
		t.Errorf("expected status 1, got %v", receivedFilter.Status)
	}
}

func TestSupplierService_List_WithMultipleFilters(t *testing.T) {
	receivedFilter := &repository.SupplierQueryFilter{}
	mockRepo := &mockSupplierRepository{
		listFunc: func(ctx context.Context, filter *repository.SupplierQueryFilter) ([]model.Supplier, int, error) {
			receivedFilter = filter
			return []model.Supplier{{BaseModel: model.BaseModel{ID: 1}}}, 1, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	status := 1

	_, err := svc.List(context.Background(), &SupplierQueryFilter{
		Code:     "SUP",
		Name:     "Test",
		Contact:  "John",
		Phone:    "123",
		Status:   &status,
		Page:     1,
		PageSize: 10,
	})

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedFilter.Code != "SUP" {
		t.Errorf("expected code 'SUP', got '%s'", receivedFilter.Code)
	}
	if receivedFilter.Name != "Test" {
		t.Errorf("expected name 'Test', got '%s'", receivedFilter.Name)
	}
	if receivedFilter.Contact != "John" {
		t.Errorf("expected contact 'John', got '%s'", receivedFilter.Contact)
	}
	if receivedFilter.Phone != "123" {
		t.Errorf("expected phone '123', got '%s'", receivedFilter.Phone)
	}
	if receivedFilter.Status == nil || *receivedFilter.Status != 1 {
		t.Errorf("expected status 1, got %v", receivedFilter.Status)
	}
}

func TestSupplierService_Update_Success(t *testing.T) {
	updatedSupplier := &model.Supplier{}
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return &model.Supplier{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Old Name",
				Code:      "SUP001",
			}, nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, supplier *model.Supplier) error {
			updatedSupplier = supplier
			return nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &UpdateSupplierInput{
		Name: strPtrSupplier("New Name"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedSupplier.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", updatedSupplier.Name)
	}
}

func TestSupplierService_Update_Code(t *testing.T) {
	updatedSupplier := &model.Supplier{}
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return &model.Supplier{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Supplier",
				Code:      "OLD",
			}, nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, supplier *model.Supplier) error {
			updatedSupplier = supplier
			return nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &UpdateSupplierInput{
		Code: strPtrSupplier("NEW"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedSupplier.Code != "NEW" {
		t.Errorf("expected code 'NEW', got '%s'", updatedSupplier.Code)
	}
}

func TestSupplierService_Update_DuplicateCode(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return &model.Supplier{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Supplier",
				Code:      "OLD",
			}, nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Supplier, error) {
			return &model.Supplier{BaseModel: model.BaseModel{ID: 2}, Code: code}, nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &UpdateSupplierInput{
		Code: strPtrSupplier("EXISTING"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestSupplierService_Update_NotFound(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewSupplierService(mockRepo, nil)
	input := &UpdateSupplierInput{Name: strPtrSupplier("Updated")}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent supplier, got nil")
	}
}

func TestSupplierService_Delete_Success(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return &model.Supplier{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestSupplierService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockSupplierRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Supplier, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewSupplierService(mockRepo, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent supplier, got nil")
	}
}

func strPtrSupplier(s string) *string {
	return &s
}
