package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockProductRepository struct {
	createFunc    func(ctx context.Context, product *model.Product) error
	getByIDFunc   func(ctx context.Context, id int64) (*model.Product, error)
	getBySKUFunc  func(ctx context.Context, sku string) (*model.Product, error)
	listFunc      func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error)
	updateFunc    func(ctx context.Context, product *model.Product) error
	deleteFunc    func(ctx context.Context, id int64) error
}

func (m *mockProductRepository) Create(ctx context.Context, product *model.Product) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, product)
	}
	return errors.New("not implemented")
}

func (m *mockProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProductRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	if m.getBySKUFunc != nil {
		return m.getBySKUFunc(ctx, sku)
	}
	return nil, errors.New("not implemented")
}

func (m *mockProductRepository) List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, categoryID, keyword)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockProductRepository) Update(ctx context.Context, product *model.Product) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, product)
	}
	return errors.New("not implemented")
}

func (m *mockProductRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestProductService_Create_Success(t *testing.T) {
	createdProduct := &model.Product{}
	mockRepo := &mockProductRepository{
		getBySKUFunc: func(ctx context.Context, sku string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
		createFunc: func(ctx context.Context, product *model.Product) error {
			product.ID = 1
			createdProduct = product
			return nil
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &CreateProductInput{
		SKU:  "SKU001",
		Name: "Test Product",
		Price: 99.99,
	}

	product, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if product == nil {
		t.Fatal("expected product, got nil")
	}
	if createdProduct.SKU != "SKU001" {
		t.Errorf("expected SKU 'SKU001', got '%s'", createdProduct.SKU)
	}
	if createdProduct.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got '%s'", createdProduct.Name)
	}
}

func TestProductService_Create_EmptySKU(t *testing.T) {
	mockRepo := &mockProductRepository{}

	svc := NewProductService(mockRepo, nil)
	input := &CreateProductInput{
		SKU:  "",
		Name: "Test Product",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty SKU, got nil")
	}
}

func TestProductService_Create_EmptyName(t *testing.T) {
	mockRepo := &mockProductRepository{}

	svc := NewProductService(mockRepo, nil)
	input := &CreateProductInput{
		SKU:  "SKU001",
		Name: "",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty name, got nil")
	}
}

func TestProductService_Create_DuplicateSKU(t *testing.T) {
	mockRepo := &mockProductRepository{
		getBySKUFunc: func(ctx context.Context, sku string) (*model.Product, error) {
			return &model.Product{SKU: sku}, nil
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &CreateProductInput{
		SKU:  "SKU001",
		Name: "Test Product",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate SKU, got nil")
	}
}

func TestProductService_Create_DefaultStatus(t *testing.T) {
	mockRepo := &mockProductRepository{
		getBySKUFunc: func(ctx context.Context, sku string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
		createFunc: func(ctx context.Context, product *model.Product) error {
			return nil
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &CreateProductInput{
		SKU:  "SKU001",
		Name: "Test Product",
	}

	product, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if product.Status != model.ProductStatusActive {
		t.Errorf("expected status %d, got %d", model.ProductStatusActive, product.Status)
	}
}

func TestProductService_GetByID_Success(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return &model.Product{
				BaseModel: model.BaseModel{ID: id},
				SKU:       "SKU001",
				Name:      "Test Product",
			}, nil
		},
	}

	svc := NewProductService(mockRepo, nil)

	product, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if product == nil {
		t.Fatal("expected product, got nil")
	}
	if product.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got '%s'", product.Name)
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewProductService(mockRepo, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent product, got nil")
	}
}

func TestProductService_List_Success(t *testing.T) {
	mockRepo := &mockProductRepository{
		listFunc: func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error) {
			return []model.Product{
				{BaseModel: model.BaseModel{ID: 1}, SKU: "SKU001", Name: "Product 1"},
				{BaseModel: model.BaseModel{ID: 2}, SKU: "SKU002", Name: "Product 2"},
			}, 2, nil
		},
	}

	svc := NewProductService(mockRepo, nil)

	result, err := svc.List(context.Background(), 1, 10, 0, "")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(result.Products))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestProductService_List_WithCategoryID(t *testing.T) {
	mockRepo := &mockProductRepository{
		listFunc: func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error) {
			if categoryID != 5 {
				t.Errorf("expected categoryID 5, got %d", categoryID)
			}
			return []model.Product{
				{BaseModel: model.BaseModel{ID: 1}, CategoryID: categoryID},
			}, 1, nil
		},
	}

	svc := NewProductService(mockRepo, nil)

	_, err := svc.List(context.Background(), 1, 10, 5, "")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestProductService_List_WithKeyword(t *testing.T) {
	mockRepo := &mockProductRepository{
		listFunc: func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error) {
			if keyword != "test" {
				t.Errorf("expected keyword 'test', got '%s'", keyword)
			}
			return []model.Product{
				{BaseModel: model.BaseModel{ID: 1}, Name: "test product"},
			}, 1, nil
		},
	}

	svc := NewProductService(mockRepo, nil)

	_, err := svc.List(context.Background(), 1, 10, 0, "test")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestProductService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockProductRepository{
		listFunc: func(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return []model.Product{}, 0, nil
		},
	}

	svc := NewProductService(mockRepo, nil)

	_, err := svc.List(context.Background(), 0, 0, 0, "")

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestProductService_Update_Success(t *testing.T) {
	updatedProduct := &model.Product{}
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return &model.Product{
				BaseModel: model.BaseModel{ID: id},
				SKU:       "SKU001",
				Name:      "Test Product",
			}, nil
		},
		getBySKUFunc: func(ctx context.Context, sku string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, product *model.Product) error {
			updatedProduct = product
			return nil
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &UpdateProductInput{
		Name:  productStrPtr("Updated Product"),
		Price: productFloatPtr(199.99),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedProduct.Name != "Updated Product" {
		t.Errorf("expected name 'Updated Product', got '%s'", updatedProduct.Name)
	}
	if updatedProduct.Price != 199.99 {
		t.Errorf("expected price 199.99, got %f", updatedProduct.Price)
	}
}

func TestProductService_Update_DuplicateSKU(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return &model.Product{
				BaseModel: model.BaseModel{ID: id},
				SKU:       "SKU001",
				Name:      "Test Product",
			}, nil
		},
		getBySKUFunc: func(ctx context.Context, sku string) (*model.Product, error) {
			if sku == "SKU002" {
				return &model.Product{BaseModel: model.BaseModel{ID: 2}, SKU: sku}, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &UpdateProductInput{
		SKU: productStrPtr("SKU002"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err == nil {
		t.Error("expected error for duplicate SKU, got nil")
	}
}

func TestProductService_Update_SameSKU(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return &model.Product{
				BaseModel: model.BaseModel{ID: id},
				SKU:       "SKU001",
				Name:      "Test Product",
			}, nil
		},
		getBySKUFunc: func(ctx context.Context, sku string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, product *model.Product) error {
			return nil
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &UpdateProductInput{
		SKU: productStrPtr("SKU001"),
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update with same SKU should succeed: %v", err)
	}
}

func TestProductService_Update_NotFound(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewProductService(mockRepo, nil)
	input := &UpdateProductInput{Name: productStrPtr("Updated")}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent product, got nil")
	}
}

func TestProductService_Delete_Success(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return &model.Product{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewProductService(mockRepo, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestProductService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockProductRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Product, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewProductService(mockRepo, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent product, got nil")
	}
}

func productStrPtr(s string) *string {
	return &s
}

func productFloatPtr(f float64) *float64 {
	return &f
}
