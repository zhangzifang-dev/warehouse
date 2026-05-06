package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"
)

type mockUserService struct {
	createFunc       func(ctx context.Context, input *service.CreateUserInput) (*model.User, error)
	getByIDFunc      func(ctx context.Context, id int64) (*model.User, error)
	listFunc         func(ctx context.Context, page, pageSize int) (*service.ListUsersResult, error)
	updateFunc       func(ctx context.Context, id int64, input *service.UpdateUserInput) (*model.User, error)
	deleteFunc       func(ctx context.Context, id int64) error
	getUserRolesFunc func(ctx context.Context, userID int64) ([]model.Role, error)
	assignRolesFunc  func(ctx context.Context, userID int64, roleIDs []int64) error
}

func (m *mockUserService) Create(ctx context.Context, input *service.CreateUserInput) (*model.User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) List(ctx context.Context, page, pageSize int) (*service.ListUsersResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) Update(ctx context.Context, id int64, input *service.UpdateUserInput) (*model.User, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockUserService) GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	if m.getUserRolesFunc != nil {
		return m.getUserRolesFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if m.assignRolesFunc != nil {
		return m.assignRolesFunc(ctx, userID, roleIDs)
	}
	return errors.New("not implemented")
}

func TestUserHandler_Create_Success(t *testing.T) {
	mockSvc := &mockUserService{
		createFunc: func(ctx context.Context, input *service.CreateUserInput) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: 1},
				Username:  input.Username,
			}, nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/users", handler.Create)

	body := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	if data["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got '%v'", data["username"])
	}
}

func TestUserHandler_Create_DuplicateUsername(t *testing.T) {
	mockSvc := &mockUserService{
		createFunc: func(ctx context.Context, input *service.CreateUserInput) (*model.User, error) {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "username already exists")
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/users", handler.Create)

	body := map[string]string{"username": "existinguser", "password": "password123"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for duplicate, got %d", w.Code)
	}
}

func TestUserHandler_Create_BindingError(t *testing.T) {
	mockSvc := &mockUserService{}
	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/users", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockUserService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
			}, nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	if data["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got '%v'", data["username"])
	}
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	mockSvc := &mockUserService{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_GetByID_InvalidID(t *testing.T) {
	mockSvc := &mockUserService{}
	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_List_Success(t *testing.T) {
	mockSvc := &mockUserService{
		listFunc: func(ctx context.Context, page, pageSize int) (*service.ListUsersResult, error) {
			return &service.ListUsersResult{
				Users: []model.User{
					{BaseModel: model.BaseModel{ID: 1}, Username: "user1"},
					{BaseModel: model.BaseModel{ID: 2}, Username: "user2"},
				},
				Total: 2,
			}, nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	users, ok := data["users"].([]interface{})
	if !ok {
		t.Fatal("expected users array in response")
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserHandler_List_DefaultPagination(t *testing.T) {
	mockSvc := &mockUserService{
		listFunc: func(ctx context.Context, page, pageSize int) (*service.ListUsersResult, error) {
			if page != 1 {
				t.Errorf("expected page 1, got %d", page)
			}
			if pageSize != 10 {
				t.Errorf("expected pageSize 10, got %d", pageSize)
			}
			return &service.ListUsersResult{Users: []model.User{}, Total: 0}, nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_Update_Success(t *testing.T) {
	mockSvc := &mockUserService{
		updateFunc: func(ctx context.Context, id int64, input *service.UpdateUserInput) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
			}, nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/users/:id", handler.Update)

	body := map[string]interface{}{"status": 1}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	mockSvc := &mockUserService{
		updateFunc: func(ctx context.Context, id int64, input *service.UpdateUserInput) (*model.User, error) {
			return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/users/:id", handler.Update)

	body := map[string]interface{}{"status": 1}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockUserService{
		deleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	mockSvc := &mockUserService{
		deleteFunc: func(ctx context.Context, id int64) error {
			return apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_GetRoles_Success(t *testing.T) {
	mockSvc := &mockUserService{
		getUserRolesFunc: func(ctx context.Context, userID int64) ([]model.Role, error) {
			return []model.Role{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "User", Code: "user"},
			}, nil
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users/:id/roles", handler.GetRoles)

	req := httptest.NewRequest(http.MethodGet, "/users/1/roles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatal("expected data array in response")
	}
	if len(data) != 2 {
		t.Errorf("expected 2 roles, got %d", len(data))
	}
}

func TestUserHandler_GetRoles_UserNotFound(t *testing.T) {
	mockSvc := &mockUserService{
		getUserRolesFunc: func(ctx context.Context, userID int64) ([]model.Role, error) {
			return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewUserHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/users/:id/roles", handler.GetRoles)

	req := httptest.NewRequest(http.MethodGet, "/users/999/roles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
