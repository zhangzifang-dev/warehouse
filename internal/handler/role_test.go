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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRoleService struct {
	listFunc   func(ctx context.Context, page, pageSize int) ([]model.Role, int, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Role, error)
	createFunc func(ctx context.Context, role *model.Role) (*model.Role, error)
	updateFunc func(ctx context.Context, id int64, role *model.Role) (*model.Role, error)
	deleteFunc func(ctx context.Context, id int64) error
}

func (m *mockRoleService) List(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockRoleService) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoleService) Create(ctx context.Context, role *model.Role) (*model.Role, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, role)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoleService) Update(ctx context.Context, id int64, role *model.Role) (*model.Role, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, role)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoleService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupRoleTest(t *testing.T) (*gin.Engine, *RoleHandler, *mockRoleService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockRoleService{}
	handler := NewRoleHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestRoleHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		mockRoles  []model.Role
		mockTotal  int
		mockError  error
		queryPage  string
		querySize  string
		wantStatus int
		wantTotal  int
	}{
		{
			name: "success with default pagination",
			mockRoles: []model.Role{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "User", Code: "user"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with custom pagination",
			mockRoles: []model.Role{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
			},
			mockTotal:  10,
			queryPage:  "2",
			querySize:  "5",
			wantStatus: http.StatusOK,
			wantTotal:  10,
		},
		{
			name:       "empty list",
			mockRoles:  []model.Role{},
			mockTotal:  0,
			wantStatus: http.StatusOK,
			wantTotal:  0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupRoleTest(t)
			mockSvc.listFunc = func(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
				if tt.mockError != nil {
					return nil, 0, tt.mockError
				}
				return tt.mockRoles, tt.mockTotal, nil
			}

			router.GET("/roles", handler.List)

			req := httptest.NewRequest("GET", "/roles?page="+tt.queryPage+"&size="+tt.querySize, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.wantTotal), data["total"])
			}
		})
	}
}

func TestRoleHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		roleID     string
		mockRole   *model.Role
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			roleID:     "1",
			mockRole:   &model.Role{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			roleID:     "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			roleID:     "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "role not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupRoleTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Role, error) {
				return tt.mockRole, tt.mockError
			}

			router.GET("/roles/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/roles/"+tt.roleID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRoleHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		mockRole   *model.Role
		mockError  error
		wantStatus int
	}{
		{
			name: "success",
			body: CreateRoleRequest{
				Name:        "Admin",
				Code:        "admin",
				Description: "Administrator role",
				Status:      1,
			},
			mockRole:   &model.Role{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateRoleRequest{Name: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate code",
			body: CreateRoleRequest{
				Name: "Admin",
				Code: "admin",
			},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "role code already exists"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupRoleTest(t)
			mockSvc.createFunc = func(ctx context.Context, role *model.Role) (*model.Role, error) {
				return tt.mockRole, tt.mockError
			}

			router.POST("/roles", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/roles", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRoleHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		roleID     string
		body       interface{}
		mockRole   *model.Role
		mockError  error
		wantStatus int
	}{
		{
			name:   "success",
			roleID: "1",
			body: UpdateRoleRequest{
				Name:        "Super Admin",
				Description: "Super administrator",
			},
			mockRole:   &model.Role{BaseModel: model.BaseModel{ID: 1}, Name: "Super Admin"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			roleID:     "invalid",
			body:       UpdateRoleRequest{Name: "Test"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			roleID:     "999",
			body:       UpdateRoleRequest{Name: "Test"},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "role not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			roleID:     "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupRoleTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, role *model.Role) (*model.Role, error) {
				return tt.mockRole, tt.mockError
			}

			router.PUT("/roles/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/roles/"+tt.roleID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRoleHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		roleID     string
		mockError  error
		wantStatus int
	}{
		{
			name:       "success",
			roleID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			roleID:     "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			roleID:     "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "role not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupRoleTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/roles/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/roles/"+tt.roleID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
