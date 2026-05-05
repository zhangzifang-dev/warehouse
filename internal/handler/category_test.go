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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockCategoryService struct {
	listFunc    func(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Category, error)
	createFunc  func(ctx context.Context, input *service.CreateCategoryInput) (*model.Category, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateCategoryInput) (*model.Category, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockCategoryService) List(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCategoryService) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCategoryService) Create(ctx context.Context, input *service.CreateCategoryInput) (*model.Category, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCategoryService) Update(ctx context.Context, id int64, input *service.UpdateCategoryInput) (*model.Category, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCategoryService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupCategoryHandlerTest(t *testing.T) (*gin.Engine, *CategoryHandler, *mockCategoryService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockCategoryService{}
	handler := NewCategoryHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestCategoryHandler_List(t *testing.T) {
	tests := []struct {
		name          string
		mockCategories []model.Category
		mockTotal     int
		mockError     error
		queryParent   string
		queryPage     string
		querySize     string
		queryName     string
		wantStatus    int
		wantTotal     int
	}{
		{
			name: "success with default pagination",
			mockCategories: []model.Category{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
				{BaseModel: model.BaseModel{ID: 2}, Name: "Clothing"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with parent filter",
			mockCategories: []model.Category{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Laptops", ParentID: 1},
			},
			mockTotal:   1,
			queryParent: "1",
			wantStatus:  http.StatusOK,
			wantTotal:   1,
		},
		{
			name: "success with name filter",
			mockCategories: []model.Category{
				{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
			},
			mockTotal:   1,
			queryName:   "Electronics",
			wantStatus:  http.StatusOK,
			wantTotal:   1,
		},
		{
			name:          "empty list",
			mockCategories: []model.Category{},
			mockTotal:     0,
			wantStatus:    http.StatusOK,
			wantTotal:     0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCategoryHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListCategoriesResult{
					Categories: tt.mockCategories,
					Total:      tt.mockTotal,
				}, nil
			}

			router.GET("/categories", handler.List)

			url := "/categories?parent_id=" + tt.queryParent + "&page=" + tt.queryPage + "&size=" + tt.querySize + "&name=" + tt.queryName
			req := httptest.NewRequest("GET", url, nil)
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

func TestCategoryHandler_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		categoryID   string
		mockCategory *model.Category
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			categoryID: "1",
			mockCategory: &model.Category{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			categoryID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			categoryID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "category not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCategoryHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Category, error) {
				return tt.mockCategory, tt.mockError
			}

			router.GET("/categories/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/categories/"+tt.categoryID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCategoryHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		mockCategory *model.Category
		mockError    error
		wantStatus   int
	}{
		{
			name: "success",
			body: CreateCategoryRequest{
				Name:      "Electronics",
				ParentID:  int64Ptr(0),
				SortOrder: intPtr(1),
			},
			mockCategory: &model.Category{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateCategoryRequest{},
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
			router, handler, mockSvc := setupCategoryHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateCategoryInput) (*model.Category, error) {
				return tt.mockCategory, tt.mockError
			}

			router.POST("/categories", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/categories", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCategoryHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		categoryID   string
		body         interface{}
		mockCategory *model.Category
		mockError    error
		wantStatus   int
	}{
		{
			name:       "success",
			categoryID: "1",
			body: UpdateCategoryRequest{
				Name:      strPtrHandler("Computers"),
				SortOrder: intPtr(2),
			},
			mockCategory: &model.Category{BaseModel: model.BaseModel{ID: 1}, Name: "Computers"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			categoryID: "invalid",
			body:       UpdateCategoryRequest{Name: strPtrHandler("Updated")},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			categoryID: "999",
			body:       UpdateCategoryRequest{Name: strPtrHandler("Updated")},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "category not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			categoryID: "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCategoryHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateCategoryInput) (*model.Category, error) {
				return tt.mockCategory, tt.mockError
			}

			router.PUT("/categories/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/categories/"+tt.categoryID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCategoryHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		categoryID  string
		mockError   error
		wantStatus  int
	}{
		{
			name:       "success",
			categoryID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			categoryID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			categoryID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "category not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "has children",
			categoryID: "1",
			mockError:  apperrors.NewAppError(apperrors.CodeBadRequest, "cannot delete category with children"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupCategoryHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/categories/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/categories/"+tt.categoryID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func int64Ptr(i int) *int64 {
	v := int64(i)
	return &v
}

func intPtr(i int) *int {
	return &i
}

func strPtrHandler(s string) *string {
	return &s
}

func TestCategoryHandler_List_WithNameFilter(t *testing.T) {
	var capturedFilter *service.CategoryQueryFilter
	mockSvc := &mockCategoryService{
		listFunc: func(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error) {
			capturedFilter = filter
			return &service.ListCategoriesResult{
				Categories: []model.Category{
					{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
				},
				Total: 1,
			}, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	router, _, _ := setupCategoryHandlerTest(t)
	router.GET("/categories", handler.List)

	req := httptest.NewRequest("GET", "/categories?name=Electronics&page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedFilter)
	assert.Equal(t, "Electronics", capturedFilter.Name)
	assert.Equal(t, 1, capturedFilter.Page)
	assert.Equal(t, 10, capturedFilter.PageSize)
}

func TestCategoryHandler_List_WithSpecialCharacters(t *testing.T) {
	var capturedFilter *service.CategoryQueryFilter
	mockSvc := &mockCategoryService{
		listFunc: func(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error) {
			capturedFilter = filter
			return &service.ListCategoriesResult{
				Categories: []model.Category{},
				Total:      0,
			}, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	router, _, _ := setupCategoryHandlerTest(t)
	router.GET("/categories", handler.List)

	req := httptest.NewRequest("GET", "/categories?name=Test%20Category%26Special%2FChars&page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedFilter)
	assert.Equal(t, "Test Category&Special/Chars", capturedFilter.Name)
}

func TestCategoryHandler_List_WithEmptyFilter(t *testing.T) {
	var capturedFilter *service.CategoryQueryFilter
	mockSvc := &mockCategoryService{
		listFunc: func(ctx context.Context, filter *service.CategoryQueryFilter) (*service.ListCategoriesResult, error) {
			capturedFilter = filter
			return &service.ListCategoriesResult{
				Categories: []model.Category{
					{BaseModel: model.BaseModel{ID: 1}, Name: "Electronics"},
					{BaseModel: model.BaseModel{ID: 2}, Name: "Clothing"},
				},
				Total: 2,
			}, nil
		},
	}

	handler := NewCategoryHandler(mockSvc)
	router, _, _ := setupCategoryHandlerTest(t)
	router.GET("/categories", handler.List)

	req := httptest.NewRequest("GET", "/categories?page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedFilter)
	assert.Equal(t, "", capturedFilter.Name)
	assert.Equal(t, 1, capturedFilter.Page)
	assert.Equal(t, 10, capturedFilter.PageSize)
}
