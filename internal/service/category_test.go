package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockCategoryRepository struct {
	createFunc      func(ctx context.Context, category *model.Category) error
	getByIDFunc     func(ctx context.Context, id int64) (*model.Category, error)
	listFunc        func(ctx context.Context, page, pageSize int, parentID int64) ([]model.Category, int, error)
	updateFunc      func(ctx context.Context, category *model.Category) error
	deleteFunc      func(ctx context.Context, id int64) error
	hasChildrenFunc func(ctx context.Context, id int64) (bool, error)
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, category)
	}
	return errors.New("not implemented")
}

func (m *mockCategoryRepository) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCategoryRepository) List(ctx context.Context, page, pageSize int, parentID int64) ([]model.Category, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, parentID)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, category)
	}
	return errors.New("not implemented")
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockCategoryRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	if m.hasChildrenFunc != nil {
		return m.hasChildrenFunc(ctx, id)
	}
	return false, errors.New("not implemented")
}

func TestCategoryService_Create_Success(t *testing.T) {
	createdCategory := &model.Category{}
	mockRepo := &mockCategoryRepository{
		createFunc: func(ctx context.Context, category *model.Category) error {
			category.ID = 1
			createdCategory = category
			return nil
		},
	}

	svc := NewCategoryService(mockRepo)
	input := &CreateCategoryInput{
		Name:      "Electronics",
		ParentID:  0,
		SortOrder: 1,
	}

	category, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if category == nil {
		t.Fatal("expected category, got nil")
	}
	if createdCategory.Name != "Electronics" {
		t.Errorf("expected name 'Electronics', got '%s'", createdCategory.Name)
	}
}

func TestCategoryService_Create_EmptyName(t *testing.T) {
	mockRepo := &mockCategoryRepository{}

	svc := NewCategoryService(mockRepo)
	input := &CreateCategoryInput{
		Name: "",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for empty name, got nil")
	}
}

func TestCategoryService_Create_DefaultStatus(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		createFunc: func(ctx context.Context, category *model.Category) error {
			return nil
		},
	}

	svc := NewCategoryService(mockRepo)
	input := &CreateCategoryInput{
		Name: "Electronics",
	}

	category, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if category.Status != model.CategoryStatusActive {
		t.Errorf("expected status %d, got %d", model.CategoryStatusActive, category.Status)
	}
}

func TestCategoryService_GetByID_Success(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return &model.Category{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Electronics",
			}, nil
		},
	}

	svc := NewCategoryService(mockRepo)

	category, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if category == nil {
		t.Fatal("expected category, got nil")
	}
	if category.Name != "Electronics" {
		t.Errorf("expected name 'Electronics', got '%s'", category.Name)
	}
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCategoryService(mockRepo)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent category, got nil")
	}
}

func TestCategoryService_List_Success(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		listFunc: func(ctx context.Context, page, pageSize int, parentID int64) ([]model.Category, int, error) {
			return []model.Category{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Clothing"},
			}, 2, nil
		},
	}

	svc := NewCategoryService(mockRepo)

	result, err := svc.List(context.Background(), 1, 10, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(result.Categories))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestCategoryService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		listFunc: func(ctx context.Context, page, pageSize int, parentID int64) ([]model.Category, int, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return []model.Category{}, 0, nil
		},
	}

	svc := NewCategoryService(mockRepo)

	_, err := svc.List(context.Background(), 0, 0, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestCategoryService_Update_Success(t *testing.T) {
	updatedCategory := &model.Category{}
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return &model.Category{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Electronics",
				SortOrder: 1,
			}, nil
		},
		updateFunc: func(ctx context.Context, category *model.Category) error {
			updatedCategory = category
			return nil
		},
	}

	svc := NewCategoryService(mockRepo)
	newSortOrder := 2
	input := &UpdateCategoryInput{
		Name:      strPtr("Computers"),
		SortOrder: &newSortOrder,
	}

	_, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedCategory.Name != "Computers" {
		t.Errorf("expected name 'Computers', got '%s'", updatedCategory.Name)
	}
	if updatedCategory.SortOrder != 2 {
		t.Errorf("expected sortOrder 2, got %d", updatedCategory.SortOrder)
	}
}

func TestCategoryService_Update_Status(t *testing.T) {
	newStatus := model.CategoryStatusInactive
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return &model.Category{
				BaseModel: model.BaseModel{ID: id},
				Name:      "Electronics",
				Status:    model.CategoryStatusActive,
			}, nil
		},
		updateFunc: func(ctx context.Context, category *model.Category) error {
			return nil
		},
	}

	svc := NewCategoryService(mockRepo)
	input := &UpdateCategoryInput{
		Status: &newStatus,
	}

	category, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if category.Status != model.CategoryStatusInactive {
		t.Errorf("expected status %d, got %d", model.CategoryStatusInactive, category.Status)
	}
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCategoryService(mockRepo)
	input := &UpdateCategoryInput{Name: strPtr("Updated")}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent category, got nil")
	}
}

func TestCategoryService_Delete_Success(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return &model.Category{BaseModel: model.BaseModel{ID: id}}, nil
		},
		hasChildrenFunc: func(ctx context.Context, id int64) (bool, error) {
			return false, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewCategoryService(mockRepo)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestCategoryService_Delete_HasChildren(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return &model.Category{BaseModel: model.BaseModel{ID: id}}, nil
		},
		hasChildrenFunc: func(ctx context.Context, id int64) (bool, error) {
			return true, nil
		},
	}

	svc := NewCategoryService(mockRepo)

	err := svc.Delete(context.Background(), 1)

	if err == nil {
		t.Error("expected error when category has children, got nil")
	}
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockCategoryRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewCategoryService(mockRepo)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent category, got nil")
	}
}

func strPtr(s string) *string {
	return &s
}
