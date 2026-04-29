package service

import (
	"context"
	"errors"
	"testing"

	"warehouse/internal/model"
)

type mockUserRepositoryForService struct {
	createFunc       func(ctx context.Context, user *model.User) error
	getByIDFunc      func(ctx context.Context, id int64) (*model.User, error)
	listFunc         func(ctx context.Context, page, pageSize int) ([]model.User, int, error)
	updateFunc       func(ctx context.Context, user *model.User) error
	deleteFunc       func(ctx context.Context, id int64) error
	getUserRolesFunc func(ctx context.Context, userID int64) ([]model.Role, error)
	assignRolesFunc  func(ctx context.Context, userID int64, roleIDs []int64) error
}

func (m *mockUserRepositoryForService) Create(ctx context.Context, user *model.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *mockUserRepositoryForService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepositoryForService) List(ctx context.Context, page, pageSize int) ([]model.User, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockUserRepositoryForService) Update(ctx context.Context, user *model.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *mockUserRepositoryForService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockUserRepositoryForService) GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	if m.getUserRolesFunc != nil {
		return m.getUserRolesFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepositoryForService) AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if m.assignRolesFunc != nil {
		return m.assignRolesFunc(ctx, userID, roleIDs)
	}
	return errors.New("not implemented")
}

func TestUserService_Create_Success(t *testing.T) {
	createdUser := &model.User{}
	mockRepo := &mockUserRepositoryForService{
		createFunc: func(ctx context.Context, user *model.User) error {
			user.ID = 1
			createdUser = user
			return nil
		},
	}

	svc := NewUserService(mockRepo, nil)
	input := &CreateUserInput{
		Username: "testuser",
		Password: "password123",
		Nickname: "Test User",
		Email:    "test@example.com",
		Phone:    "1234567890",
	}

	user, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if createdUser.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", createdUser.Username)
	}
	if createdUser.Password == "password123" {
		t.Error("expected password to be hashed")
	}
}

func TestUserService_Create_DefaultStatus(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		createFunc: func(ctx context.Context, user *model.User) error {
			return nil
		},
	}

	svc := NewUserService(mockRepo, nil)
	input := &CreateUserInput{
		Username: "testuser",
		Password: "password123",
	}

	user, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user.Status != model.UserStatusActive {
		t.Errorf("expected status %d, got %d", model.UserStatusActive, user.Status)
	}
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		createFunc: func(ctx context.Context, user *model.User) error {
			return errors.New("duplicate key")
		},
	}

	svc := NewUserService(mockRepo, nil)
	input := &CreateUserInput{
		Username: "existinguser",
		Password: "password123",
	}

	_, err := svc.Create(context.Background(), input)

	if err == nil {
		t.Error("expected error for duplicate username, got nil")
	}
}

func TestUserService_GetByID_Success(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Nickname:  "Test User",
			}, nil
		},
	}

	svc := NewUserService(mockRepo, nil)

	user, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewUserService(mockRepo, nil)

	_, err := svc.GetByID(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestUserService_List_Success(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.User, int, error) {
			return []model.User{
				{BaseModel: model.BaseModel{ID: 1}, Username: "user1"},
				{BaseModel: model.BaseModel{ID: 2}, Username: "user2"},
			}, 2, nil
		},
	}

	svc := NewUserService(mockRepo, nil)

	result, err := svc.List(context.Background(), 1, 10)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(result.Users))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestUserService_List_DefaultPagination(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.User, int, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return []model.User{}, 0, nil
		},
	}

	svc := NewUserService(mockRepo, nil)

	_, err := svc.List(context.Background(), 0, 0)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestUserService_List_MaxPageSize(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		listFunc: func(ctx context.Context, page, pageSize int) ([]model.User, int, error) {
			if pageSize > 100 {
				t.Errorf("expected pageSize <= 100, got %d", pageSize)
			}
			return []model.User{}, 0, nil
		},
	}

	svc := NewUserService(mockRepo, nil)

	_, err := svc.List(context.Background(), 1, 200)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
}

func TestUserService_Update_Success(t *testing.T) {
	updatedUser := &model.User{}
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Nickname:  "Old Nickname",
			}, nil
		},
		updateFunc: func(ctx context.Context, user *model.User) error {
			updatedUser = user
			return nil
		},
	}

	svc := NewUserService(mockRepo, nil)
	input := &UpdateUserInput{
		Nickname: "New Nickname",
		Email:    "new@example.com",
	}

	user, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if user.Nickname != "New Nickname" {
		t.Errorf("expected nickname 'New Nickname', got '%s'", updatedUser.Nickname)
	}
}

func TestUserService_Update_Status(t *testing.T) {
	newStatus := model.UserStatusDisabled
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Status:    model.UserStatusActive,
			}, nil
		},
		updateFunc: func(ctx context.Context, user *model.User) error {
			return nil
		},
	}

	svc := NewUserService(mockRepo, nil)
	input := &UpdateUserInput{
		Status: &newStatus,
	}

	user, err := svc.Update(context.Background(), 1, input)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if user.Status != model.UserStatusDisabled {
		t.Errorf("expected status %d, got %d", model.UserStatusDisabled, user.Status)
	}
}

func TestUserService_Update_NotFound(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewUserService(mockRepo, nil)
	input := &UpdateUserInput{Nickname: "New Name"}

	_, err := svc.Update(context.Background(), 999, input)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestUserService_Delete_Success(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{BaseModel: model.BaseModel{ID: id}}, nil
		},
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	svc := NewUserService(mockRepo, nil)

	err := svc.Delete(context.Background(), 1)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestUserService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewUserService(mockRepo, nil)

	err := svc.Delete(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestUserService_GetUserRoles_Success(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{BaseModel: model.BaseModel{ID: id}}, nil
		},
		getUserRolesFunc: func(ctx context.Context, userID int64) ([]model.Role, error) {
			return []model.Role{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "User", Code: "user"},
			}, nil
		},
	}

	svc := NewUserService(mockRepo, nil)

	roles, err := svc.GetUserRoles(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetUserRoles failed: %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
}

func TestUserService_GetUserRoles_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepositoryForService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}

	svc := NewUserService(mockRepo, nil)

	_, err := svc.GetUserRoles(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}
