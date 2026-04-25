package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockRoleRepository struct {
	createFunc            func(ctx context.Context, role *model.Role) error
	getByIDFunc           func(ctx context.Context, id int64) (*model.Role, error)
	getByCodeFunc         func(ctx context.Context, code string) (*model.Role, error)
	listFunc              func(ctx context.Context, page, pageSize int) ([]model.Role, int, error)
	updateFunc            func(ctx context.Context, role *model.Role) error
	deleteFunc            func(ctx context.Context, id int64) error
	assignPermissionsFunc func(ctx context.Context, roleID int64, permissionIDs []int64) error
	getRolePermissionsFunc func(ctx context.Context, roleID int64) ([]model.Permission, error)
}

func (m *mockRoleRepository) Create(ctx context.Context, role *model.Role) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, role)
	}
	return errors.New("not implemented")
}

func (m *mockRoleRepository) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	if m.getByCodeFunc != nil {
		return m.getByCodeFunc(ctx, code)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoleRepository) List(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockRoleRepository) Update(ctx context.Context, role *model.Role) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, role)
	}
	return errors.New("not implemented")
}

func (m *mockRoleRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRoleRepository) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if m.assignPermissionsFunc != nil {
		return m.assignPermissionsFunc(ctx, roleID, permissionIDs)
	}
	return errors.New("not implemented")
}

func (m *mockRoleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error) {
	if m.getRolePermissionsFunc != nil {
		return m.getRolePermissionsFunc(ctx, roleID)
	}
	return nil, errors.New("not implemented")
}

func TestRoleService_Create_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		createFunc: func(ctx context.Context, role *model.Role) error {
			role.ID = 1
			return nil
		},
	}

	svc := NewRoleService(mockRepo)

	role := &model.Role{
		Name:        "Admin",
		Code:        "admin",
		Description: "Administrator role",
		Status:      1,
	}

	created, err := svc.Create(context.Background(), role)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("expected role ID 1, got %d", created.ID)
	}
}

func TestRoleService_Create_DuplicateCode(t *testing.T) {
	mockRepo := &mockRoleRepository{
		createFunc: func(ctx context.Context, role *model.Role) error {
			return errors.New("duplicate key")
		},
	}

	svc := NewRoleService(mockRepo)

	role := &model.Role{
		Name: "Admin",
		Code: "admin",
	}

	_, err := svc.Create(context.Background(), role)

	if err == nil {
		t.Error("expected error for duplicate code, got nil")
	}
}

func TestRoleService_GetByID_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Role, error) {
			return &model.Role{
				BaseModel:   model.BaseModel{ID: id},
				Name:        "Admin",
				Code:        "admin",
				Description: "Administrator role",
				Status:      1,
			}, nil
		},
	}

	svc := NewRoleService(mockRepo)

	role, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if role == nil {
		t.Fatal("expected role, got nil")
	}
	if role.Code != "admin" {
		t.Errorf("expected code 'admin', got '%s'", role.Code)
	}
}

func TestRoleService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Role, error) {
			return nil, errors.New("record not found")
		},
	}

	svc := NewRoleService(mockRepo)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent role, got nil")
	}
}

func TestRoleService_GetByCode_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getByCodeFunc: func(ctx context.Context, code string) (*model.Role, error) {
			return &model.Role{
				BaseModel:   model.BaseModel{ID: 1},
				Name:        "Admin",
				Code:        code,
				Description: "Administrator role",
			}, nil
		},
	}

	svc := NewRoleService(mockRepo)

	role, err := svc.GetByCode(context.Background(), "admin")

	if err != nil {
		t.Fatalf("GetByCode failed: %v", err)
	}
	if role == nil {
		t.Fatal("expected role, got nil")
	}
	if role.Name != "Admin" {
		t.Errorf("expected name 'Admin', got '%s'", role.Name)
	}
}

func TestRoleService_GetByCode_NotFound(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getByCodeFunc: func(ctx context.Context, code string) (*model.Role, error) {
			return nil, errors.New("record not found")
		},
	}

	svc := NewRoleService(mockRepo)

	_, err := svc.GetByCode(context.Background(), "nonexistent")

	if err == nil {
		t.Error("expected error for non-existent code, got nil")
	}
}

func TestRoleService_List_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
			return []model.Role{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "User", Code: "user"},
			}, 2, nil
		},
	}

	svc := NewRoleService(mockRepo)

	roles, total, err := svc.List(context.Background(), 1, 10)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestRoleService_Update_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Role, error) {
			return &model.Role{
				BaseModel:   model.BaseModel{ID: id},
				Name:        "Admin",
				Code:        "admin",
				Description: "Administrator role",
				Status:      1,
			}, nil
		},
		updateFunc: func(ctx context.Context, role *model.Role) error {
			return nil
		},
	}

	svc := NewRoleService(mockRepo)

	role := &model.Role{
		Name:        "Admin Updated",
		Description: "Updated description",
	}

	updated, err := svc.Update(context.Background(), 1, role)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Admin Updated" {
		t.Errorf("expected name 'Admin Updated', got '%s'", updated.Name)
	}
}

func TestRoleService_Update_NotFound(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.Role, error) {
			return nil, errors.New("record not found")
		},
	}

	svc := NewRoleService(mockRepo)

	role := &model.Role{
		Name: "Admin",
	}

	_, err := svc.Update(context.Background(), 999, role)

	if err == nil {
		t.Error("expected error for non-existent role, got nil")
	}
}

func TestRoleService_Delete_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewRoleService(mockRepo)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestRoleService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockRoleRepository{
		deleteFunc: func(ctx context.Context, id int64) error {
			return errors.New("record not found")
		},
	}

	svc := NewRoleService(mockRepo)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent role, got nil")
	}
}

func TestRoleService_AssignPermissions_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		assignPermissionsFunc: func(ctx context.Context, roleID int64, permissionIDs []int64) error {
			return nil
		},
	}

	svc := NewRoleService(mockRepo)

	err := svc.AssignPermissions(context.Background(), 1, []int64{1, 2, 3})

	if err != nil {
		t.Fatalf("AssignPermissions failed: %v", err)
	}
}

func TestRoleService_GetRolePermissions_Success(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getRolePermissionsFunc: func(ctx context.Context, roleID int64) ([]model.Permission, error) {
			return []model.Permission{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Create User", Code: "user:create"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Delete User", Code: "user:delete"},
			}, nil
		},
	}

	svc := NewRoleService(mockRepo)

	permissions, err := svc.GetRolePermissions(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetRolePermissions failed: %v", err)
	}
	if len(permissions) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(permissions))
	}
}

func TestRoleService_GetRolePermissions_Empty(t *testing.T) {
	mockRepo := &mockRoleRepository{
		getRolePermissionsFunc: func(ctx context.Context, roleID int64) ([]model.Permission, error) {
			return []model.Permission{}, nil
		},
	}

	svc := NewRoleService(mockRepo)

	permissions, err := svc.GetRolePermissions(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetRolePermissions failed: %v", err)
	}
	if len(permissions) != 0 {
		t.Errorf("expected 0 permissions, got %d", len(permissions))
	}
}
