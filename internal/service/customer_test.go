package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockCustomerRepository struct {
	createFunc    func(ctx context.Context, customer *model.Customer) error
	getByIDFunc   func(ctx context.Context, id int64) (*model.Customer, error)
	getByCodeFunc func(ctx context.Context, code string) (*model.Customer, error)
	listFunc      func(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error)
	updateFunc    func(ctx context.Context, customer *model.Customer) error
	deleteFunc    func(ctx context.Context, id int64) error
}

func (m *mockCustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, customer)
	}
	return errors.New("not implemented")
}

func (m *mockCustomerRepository) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCustomerRepository) GetByCode(ctx context.Context, code string) (*model.Customer, error) {
	if m.getByCodeFunc != nil {
		return m.getByCodeFunc(ctx, code)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCustomerRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, keyword)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockCustomerRepository) Update(ctx context.Context, customer *model.Customer) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, customer)
	}
	return errors.New("not implemented")
}

func (m *mockCustomerRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestCustomerService_Create_Success(t *testing.T) {
	createdCustomer := &model.Customer{}
	mockRepo := &mockCustomerRepository{
		createFunc: func(ctx context.Context, customer *model.Customer) error {
			customer.ID = 1
			createdCustomer = customer
			return nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &CreateCustomerInput{
		Name:    "Test Customer",
		Code:    "CUS001",
		Contact: "John Doe",
		Phone:   "1234567890",
	}

	customer, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if customer == nil {
		t.Fatal("expected customer, got nil")
	}
	if createdCustomer.Name != "Test Customer" {
		t.Errorf("expected name 'Test Customer', got '%s'", createdCustomer.Name)
	}
}

func TestCustomerService_Create_EmptyName(t *testing.T) {
	mockRepo := &mockCustomerRepository{}

	svc := NewCustomerService(mockRepo, nil)
	input := &CreateCustomerInput{
		Name: "",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty name, got nil")
	}
}

func TestCustomerService_Create_DuplicateCode(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByCodeFunc: func(ctx context.Context, code string) (*model.Customer, error) {
			return &model.Customer{BaseModel: model.BaseModel{ID: 1}, Code: code}, nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &CreateCustomerInput{
		Name: "Test Customer",
		Code: "CUS001",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestCustomerService_Create_DefaultStatus(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		createFunc: func(ctx context.Context, customer *model.Customer) error {
			return nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &CreateCustomerInput{
		Name: "Test Customer",
	}

	customer, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if customer.Status != model.CustomerStatusActive {
		t.Errorf("expected status %d, got %d", model.CustomerStatusActive, customer.Status)
	}
}

func TestCustomerService_GetByID_Success(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return &model.Customer{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Customer",
			}, nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	customer, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if customer == nil {
		t.Fatal("expected customer, got nil")
	}
	if customer.Name != "Test Customer" {
		t.Errorf("expected name 'Test Customer', got '%s'", customer.Name)
	}
}

func TestCustomerService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent customer, got nil")
	}
}

func TestCustomerService_List_Success(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		listFunc: func(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error) {
			return []model.Customer{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Customer A"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Customer B"},
			}, 2, nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	result, err := svc.List(context.Background(), 1, 10, "")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(result.Customers))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestCustomerService_List_WithKeyword(t *testing.T) {
	receivedKeyword := ""
	mockRepo := &mockCustomerRepository{
		listFunc: func(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error) {
			receivedKeyword = keyword
			return []model.Customer{{BaseModel: model.BaseModel{ID: 1}, Name: "Test Customer"}}, 1, nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	_, err := svc.List(context.Background(), 1, 10, "test")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if receivedKeyword != "test" {
		t.Errorf("expected keyword 'test', got '%s'", receivedKeyword)
	}
}

func TestCustomerService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		listFunc: func(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return []model.Customer{}, 0, nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	_, err := svc.List(context.Background(), 0, 0, "")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestCustomerService_Update_Success(t *testing.T) {
	updatedCustomer := &model.Customer{}
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return &model.Customer{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Old Name",
				Code:      "CUS001",
			}, nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, customer *model.Customer) error {
			updatedCustomer = customer
			return nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &UpdateCustomerInput{
		Name: strPtrCustomer("New Name"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedCustomer.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", updatedCustomer.Name)
	}
}

func TestCustomerService_Update_Code(t *testing.T) {
	updatedCustomer := &model.Customer{}
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return &model.Customer{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Customer",
				Code:      "OLD",
			}, nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, customer *model.Customer) error {
			updatedCustomer = customer
			return nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &UpdateCustomerInput{
		Code: strPtrCustomer("NEW"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedCustomer.Code != "NEW" {
		t.Errorf("expected code 'NEW', got '%s'", updatedCustomer.Code)
	}
}

func TestCustomerService_Update_DuplicateCode(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return &model.Customer{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Test Customer",
				Code:      "OLD",
			}, nil
		},
		getByCodeFunc: func(ctx context.Context, code string) (*model.Customer, error) {
			return &model.Customer{BaseModel: model.BaseModel{ID: 2}, Code: code}, nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &UpdateCustomerInput{
		Code: strPtrCustomer("EXISTING"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestCustomerService_Update_NotFound(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCustomerService(mockRepo, nil)
	input := &UpdateCustomerInput{Name: strPtrCustomer("Updated")}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent customer, got nil")
	}
}

func TestCustomerService_Delete_Success(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return &model.Customer{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestCustomerService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockCustomerRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Customer, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCustomerService(mockRepo, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent customer, got nil")
	}
}

func strPtrCustomer(s string) *string {
	return &s
}
