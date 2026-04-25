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

type mockLocationService struct {
	listFunc    func(ctx context.Context, page, pageSize int, warehouseID int64) (*service.ListLocationsResult, error)
	getByIDFunc func(ctx context.Context, id int64) (*model.Location, error)
	createFunc  func(ctx context.Context, input *service.CreateLocationInput) (*model.Location, error)
	updateFunc  func(ctx context.Context, id int64, input *service.UpdateLocationInput) (*model.Location, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockLocationService) List(ctx context.Context, page, pageSize int, warehouseID int64) (*service.ListLocationsResult, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, page, pageSize, warehouseID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLocationService) GetByID(ctx context.Context, id int64) (*model.Location, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLocationService) Create(ctx context.Context, input *service.CreateLocationInput) (*model.Location, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLocationService) Update(ctx context.Context, id int64, input *service.UpdateLocationInput) (*model.Location, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLocationService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func setupLocationHandlerTest(t *testing.T) (*gin.Engine, *LocationHandler, *mockLocationService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockLocationService{}
	handler := NewLocationHandler(mockSvc)
	router := gin.New()
	return router, handler, mockSvc
}

func TestLocationHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		mockLocations  []model.Location
		mockTotal      int
		mockError      error
		queryWarehouse string
		queryPage      string
		querySize      string
		wantStatus     int
		wantTotal      int
	}{
		{
			name: "success with default pagination",
			mockLocations: []model.Location{
				{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, Code: "A-01-02-03"},
				{BaseModel: model.BaseModel{ID: 2}, WarehouseID: 1, Code: "A-01-02-04"},
			},
			mockTotal:  2,
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name: "success with warehouse filter",
			mockLocations: []model.Location{
				{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, Code: "A-01-02-03"},
			},
			mockTotal:      1,
			queryWarehouse: "1",
			wantStatus:     http.StatusOK,
			wantTotal:      1,
		},
		{
			name:           "empty list",
			mockLocations:  []model.Location{},
			mockTotal:      0,
			wantStatus:     http.StatusOK,
			wantTotal:      0,
		},
		{
			name:       "service error",
			mockError:  apperrors.NewAppError(apperrors.CodeInternalError, "database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupLocationHandlerTest(t)
			mockSvc.listFunc = func(ctx context.Context, page, pageSize int, warehouseID int64) (*service.ListLocationsResult, error) {
				if tt.mockError != nil {
					return nil, tt.mockError
				}
				return &service.ListLocationsResult{
					Locations: tt.mockLocations,
					Total:     tt.mockTotal,
				}, nil
			}

			router.GET("/locations", handler.List)

			req := httptest.NewRequest("GET", "/locations?warehouse_id="+tt.queryWarehouse+"&page="+tt.queryPage+"&size="+tt.querySize, nil)
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

func TestLocationHandler_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		locationID    string
		mockLocation  *model.Location
		mockError     error
		wantStatus    int
	}{
		{
			name:       "success",
			locationID: "1",
			mockLocation: &model.Location{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, Code: "A-01-02-03"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			locationID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			locationID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "location not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupLocationHandlerTest(t)
			mockSvc.getByIDFunc = func(ctx context.Context, id int64) (*model.Location, error) {
				return tt.mockLocation, tt.mockError
			}

			router.GET("/locations/:id", handler.GetByID)

			req := httptest.NewRequest("GET", "/locations/"+tt.locationID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLocationHandler_Create(t *testing.T) {
	tests := []struct {
		name          string
		body          interface{}
		mockLocation  *model.Location
		mockError     error
		wantStatus    int
	}{
		{
			name: "success",
			body: CreateLocationRequest{
				WarehouseID: 1,
				Zone:        "A",
				Shelf:       "01",
				Level:       "02",
				Position:    "03",
			},
			mockLocation: &model.Location{BaseModel: model.BaseModel{ID: 1}, WarehouseID: 1, Code: "A-01-02-03"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "missing required fields",
			body:       CreateLocationRequest{Zone: "A"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate code",
			body: CreateLocationRequest{
				WarehouseID: 1,
				Zone:        "A",
				Shelf:       "01",
				Level:       "02",
				Position:    "03",
			},
			mockError:  apperrors.NewAppError(apperrors.CodeDuplicateEntry, "location code already exists"),
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
			router, handler, mockSvc := setupLocationHandlerTest(t)
			mockSvc.createFunc = func(ctx context.Context, input *service.CreateLocationInput) (*model.Location, error) {
				return tt.mockLocation, tt.mockError
			}

			router.POST("/locations", handler.Create)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("POST", "/locations", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLocationHandler_Update(t *testing.T) {
	tests := []struct {
		name          string
		locationID    string
		body          interface{}
		mockLocation  *model.Location
		mockError     error
		wantStatus    int
	}{
		{
			name:       "success",
			locationID: "1",
			body: UpdateLocationRequest{
				Zone:  "B",
				Shelf: "02",
			},
			mockLocation: &model.Location{BaseModel: model.BaseModel{ID: 1}, Zone: "B", Code: "B-02-02-03"},
			wantStatus:   http.StatusOK,
		},
		{
			name:       "invalid id",
			locationID: "invalid",
			body:       UpdateLocationRequest{Zone: "B"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			locationID: "999",
			body:       UpdateLocationRequest{Zone: "B"},
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "location not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid json",
			locationID: "1",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupLocationHandlerTest(t)
			mockSvc.updateFunc = func(ctx context.Context, id int64, input *service.UpdateLocationInput) (*model.Location, error) {
				return tt.mockLocation, tt.mockError
			}

			router.PUT("/locations/:id", handler.Update)

			var body bytes.Buffer
			if str, ok := tt.body.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.body)
			}

			req := httptest.NewRequest("PUT", "/locations/"+tt.locationID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLocationHandler_Delete(t *testing.T) {
	tests := []struct {
		name        string
		locationID  string
		mockError   error
		wantStatus  int
	}{
		{
			name:       "success",
			locationID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			locationID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			locationID: "999",
			mockError:  apperrors.NewAppError(apperrors.CodeNotFound, "location not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler, mockSvc := setupLocationHandlerTest(t)
			mockSvc.deleteFunc = func(ctx context.Context, id int64) error {
				return tt.mockError
			}

			router.DELETE("/locations/:id", handler.Delete)

			req := httptest.NewRequest("DELETE", "/locations/"+tt.locationID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
