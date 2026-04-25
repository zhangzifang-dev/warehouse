package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"warehouse/internal/handler"
	"warehouse/internal/model"
	"warehouse/internal/pkg/jwt"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestJWTService() *jwt.JWT {
	return jwt.NewJWT("test-secret-key", time.Hour)
}

func newTestRouter() *gin.Engine {
	return gin.New()
}

func TestSetup_PublicRoutes(t *testing.T) {
	r := newTestRouter()
	jwtSvc := newTestJWTService()

	handlers := &Handlers{
		Auth: handler.NewAuthHandler(&mockAuthService{}),
	}

	Setup(r, jwtSvc, handlers)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized {
		t.Error("public route should not require auth")
	}
}

func TestSetup_ProtectedRoutes_RequireAuth(t *testing.T) {
	r := newTestRouter()
	jwtSvc := newTestJWTService()

	handlers := &Handlers{
		Auth:      handler.NewAuthHandler(&mockAuthService{}),
		User:      handler.NewUserHandler(&mockUserService{}),
		Role:      handler.NewRoleHandler(&mockRoleService{}),
		Warehouse: handler.NewWarehouseHandler(&mockWarehouseService{}),
	}

	Setup(r, jwtSvc, handlers)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"GET /auth/profile", http.MethodGet, "/api/v1/auth/profile"},
		{"PUT /auth/password", http.MethodPut, "/api/v1/auth/password"},
		{"GET /users", http.MethodGet, "/api/v1/users"},
		{"GET /roles", http.MethodGet, "/api/v1/roles"},
		{"GET /warehouses", http.MethodGet, "/api/v1/warehouses"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401 for %s, got %d", tt.path, w.Code)
			}
		})
	}
}

func TestSetup_ProtectedRoutes_WithValidToken(t *testing.T) {
	r := newTestRouter()
	jwtSvc := newTestJWTService()

	handlers := &Handlers{
		Auth:      handler.NewAuthHandler(&mockAuthService{profile: &model.User{}}),
		User:      handler.NewUserHandler(&mockUserService{}),
		Role:      handler.NewRoleHandler(&mockRoleService{}),
		Warehouse: handler.NewWarehouseHandler(&mockWarehouseService{}),
		Location:  handler.NewLocationHandler(&mockLocationService{}),
		Category:  handler.NewCategoryHandler(&mockCategoryService{}),
		Product:   handler.NewProductHandler(&mockProductService{}),
		Inventory: handler.NewInventoryHandler(&mockInventoryService{}),
		Supplier:  handler.NewSupplierHandler(&mockSupplierService{}),
		Customer:  handler.NewCustomerHandler(&mockCustomerService{}),
	}

	Setup(r, jwtSvc, handlers)

	token, _ := jwtSvc.GenerateToken(1, "testuser")

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"GET /auth/profile", http.MethodGet, "/api/v1/auth/profile"},
		{"GET /users", http.MethodGet, "/api/v1/users"},
		{"GET /roles", http.MethodGet, "/api/v1/roles"},
		{"GET /warehouses", http.MethodGet, "/api/v1/warehouses"},
		{"GET /locations", http.MethodGet, "/api/v1/locations"},
		{"GET /categories", http.MethodGet, "/api/v1/categories"},
		{"GET /products", http.MethodGet, "/api/v1/products"},
		{"GET /inventory", http.MethodGet, "/api/v1/inventory"},
		{"GET /suppliers", http.MethodGet, "/api/v1/suppliers"},
		{"GET /customers", http.MethodGet, "/api/v1/customers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code == http.StatusUnauthorized {
				t.Errorf("expected not 401 for %s with valid token", tt.path)
			}
		})
	}
}

func TestSetup_AllRoutesRegistered(t *testing.T) {
	r := newTestRouter()
	jwtSvc := newTestJWTService()

	handlers := &Handlers{
		Auth:          handler.NewAuthHandler(&mockAuthService{}),
		User:          handler.NewUserHandler(&mockUserService{}),
		Role:          handler.NewRoleHandler(&mockRoleService{}),
		Warehouse:     handler.NewWarehouseHandler(&mockWarehouseService{}),
		Location:      handler.NewLocationHandler(&mockLocationService{}),
		Category:      handler.NewCategoryHandler(&mockCategoryService{}),
		Product:       handler.NewProductHandler(&mockProductService{}),
		Inventory:     handler.NewInventoryHandler(&mockInventoryService{}),
		Supplier:      handler.NewSupplierHandler(&mockSupplierService{}),
		Customer:      handler.NewCustomerHandler(&mockCustomerService{}),
		InboundOrder:  handler.NewInboundOrderHandler(&mockInboundOrderService{}),
		OutboundOrder: handler.NewOutboundOrderHandler(&mockOutboundOrderService{}),
		StockTransfer: handler.NewStockTransferHandler(&mockStockTransferService{}),
		AuditLog:      handler.NewAuditLogHandler(&mockAuditLogService{}),
	}

	Setup(r, jwtSvc, handlers)

	routes := r.Routes()

	routeMap := make(map[string]bool)
	for _, route := range routes {
		key := route.Method + " " + route.Path
		routeMap[key] = true
	}

	expectedRoutes := []string{
		"POST /api/v1/auth/login",
		"GET /api/v1/auth/profile",
		"PUT /api/v1/auth/password",
		"GET /api/v1/users",
		"POST /api/v1/users",
		"GET /api/v1/users/:id",
		"PUT /api/v1/users/:id",
		"DELETE /api/v1/users/:id",
		"GET /api/v1/users/:id/roles",
		"GET /api/v1/roles",
		"POST /api/v1/roles",
		"GET /api/v1/roles/:id",
		"PUT /api/v1/roles/:id",
		"DELETE /api/v1/roles/:id",
		"GET /api/v1/warehouses",
		"POST /api/v1/warehouses",
		"GET /api/v1/warehouses/:id",
		"PUT /api/v1/warehouses/:id",
		"DELETE /api/v1/warehouses/:id",
		"GET /api/v1/locations",
		"POST /api/v1/locations",
		"GET /api/v1/locations/:id",
		"PUT /api/v1/locations/:id",
		"DELETE /api/v1/locations/:id",
		"GET /api/v1/categories",
		"POST /api/v1/categories",
		"GET /api/v1/categories/:id",
		"PUT /api/v1/categories/:id",
		"DELETE /api/v1/categories/:id",
		"GET /api/v1/products",
		"POST /api/v1/products",
		"GET /api/v1/products/:id",
		"PUT /api/v1/products/:id",
		"DELETE /api/v1/products/:id",
		"GET /api/v1/inventory",
		"POST /api/v1/inventory",
		"GET /api/v1/inventory/:id",
		"PUT /api/v1/inventory/:id",
		"DELETE /api/v1/inventory/:id",
		"POST /api/v1/inventory/adjust",
		"POST /api/v1/inventory/check",
		"GET /api/v1/suppliers",
		"POST /api/v1/suppliers",
		"GET /api/v1/suppliers/:id",
		"PUT /api/v1/suppliers/:id",
		"DELETE /api/v1/suppliers/:id",
		"GET /api/v1/customers",
		"POST /api/v1/customers",
		"GET /api/v1/customers/:id",
		"PUT /api/v1/customers/:id",
		"DELETE /api/v1/customers/:id",
		"GET /api/v1/inbound-orders",
		"POST /api/v1/inbound-orders",
		"GET /api/v1/inbound-orders/:id",
		"PUT /api/v1/inbound-orders/:id",
		"DELETE /api/v1/inbound-orders/:id",
		"POST /api/v1/inbound-orders/:id/confirm",
		"GET /api/v1/outbound-orders",
		"POST /api/v1/outbound-orders",
		"GET /api/v1/outbound-orders/:id",
		"PUT /api/v1/outbound-orders/:id",
		"DELETE /api/v1/outbound-orders/:id",
		"POST /api/v1/outbound-orders/:id/confirm",
		"GET /api/v1/stock-transfers",
		"POST /api/v1/stock-transfers",
		"GET /api/v1/stock-transfers/:id",
		"PUT /api/v1/stock-transfers/:id",
		"DELETE /api/v1/stock-transfers/:id",
		"POST /api/v1/stock-transfers/:id/confirm",
		"GET /api/v1/audit-logs",
		"GET /api/v1/audit-logs/:id",
	}

	for _, expected := range expectedRoutes {
		if !routeMap[expected] {
			t.Errorf("expected route %s not found", expected)
		}
	}
}

type mockAuthService struct {
	profile *model.User
}

func (m *mockAuthService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	return "", nil, nil
}

func (m *mockAuthService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	if m.profile != nil {
		return m.profile, nil
	}
	return &model.User{}, nil
}

func (m *mockAuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	return nil
}

type mockUserService struct{}

func (m *mockUserService) Create(ctx context.Context, input *service.CreateUserInput) (*model.User, error) {
	return &model.User{}, nil
}

func (m *mockUserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return &model.User{}, nil
}

func (m *mockUserService) List(ctx context.Context, page, pageSize int) (*service.ListUsersResult, error) {
	return &service.ListUsersResult{}, nil
}

func (m *mockUserService) Update(ctx context.Context, id int64, input *service.UpdateUserInput) (*model.User, error) {
	return &model.User{}, nil
}

func (m *mockUserService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockUserService) GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	return nil, nil
}

type mockRoleService struct{}

func (m *mockRoleService) List(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
	return nil, 0, nil
}

func (m *mockRoleService) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	return &model.Role{}, nil
}

func (m *mockRoleService) Create(ctx context.Context, role *model.Role) (*model.Role, error) {
	return &model.Role{}, nil
}

func (m *mockRoleService) Update(ctx context.Context, id int64, role *model.Role) (*model.Role, error) {
	return &model.Role{}, nil
}

func (m *mockRoleService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockWarehouseService struct{}

func (m *mockWarehouseService) Create(ctx context.Context, input *service.CreateWarehouseInput) (*model.Warehouse, error) {
	return &model.Warehouse{}, nil
}

func (m *mockWarehouseService) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	return &model.Warehouse{}, nil
}

func (m *mockWarehouseService) List(ctx context.Context, page, pageSize int) (*service.ListWarehousesResult, error) {
	return &service.ListWarehousesResult{}, nil
}

func (m *mockWarehouseService) Update(ctx context.Context, id int64, input *service.UpdateWarehouseInput) (*model.Warehouse, error) {
	return &model.Warehouse{}, nil
}

func (m *mockWarehouseService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockLocationService struct{}

func (m *mockLocationService) Create(ctx context.Context, input *service.CreateLocationInput) (*model.Location, error) {
	return &model.Location{}, nil
}

func (m *mockLocationService) GetByID(ctx context.Context, id int64) (*model.Location, error) {
	return &model.Location{}, nil
}

func (m *mockLocationService) List(ctx context.Context, page, pageSize int, warehouseID int64) (*service.ListLocationsResult, error) {
	return &service.ListLocationsResult{}, nil
}

func (m *mockLocationService) Update(ctx context.Context, id int64, input *service.UpdateLocationInput) (*model.Location, error) {
	return &model.Location{}, nil
}

func (m *mockLocationService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockCategoryService struct{}

func (m *mockCategoryService) Create(ctx context.Context, input *service.CreateCategoryInput) (*model.Category, error) {
	return &model.Category{}, nil
}

func (m *mockCategoryService) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	return &model.Category{}, nil
}

func (m *mockCategoryService) List(ctx context.Context, page, pageSize int, parentID int64) (*service.ListCategoriesResult, error) {
	return &service.ListCategoriesResult{}, nil
}

func (m *mockCategoryService) Update(ctx context.Context, id int64, input *service.UpdateCategoryInput) (*model.Category, error) {
	return &model.Category{}, nil
}

func (m *mockCategoryService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockProductService struct{}

func (m *mockProductService) Create(ctx context.Context, input *service.CreateProductInput) (*model.Product, error) {
	return &model.Product{}, nil
}

func (m *mockProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return &model.Product{}, nil
}

func (m *mockProductService) List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) (*service.ListProductsResult, error) {
	return &service.ListProductsResult{}, nil
}

func (m *mockProductService) Update(ctx context.Context, id int64, input *service.UpdateProductInput) (*model.Product, error) {
	return &model.Product{}, nil
}

func (m *mockProductService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockInventoryService struct{}

func (m *mockInventoryService) Create(ctx context.Context, input *service.CreateInventoryInput) (*model.Inventory, error) {
	return &model.Inventory{}, nil
}

func (m *mockInventoryService) GetByID(ctx context.Context, id int64) (*model.Inventory, error) {
	return &model.Inventory{}, nil
}

func (m *mockInventoryService) List(ctx context.Context, page, pageSize int, warehouseID, productID int64) (*service.ListInventoriesResult, error) {
	return &service.ListInventoriesResult{}, nil
}

func (m *mockInventoryService) Update(ctx context.Context, id int64, input *service.UpdateInventoryInput) (*model.Inventory, error) {
	return &model.Inventory{}, nil
}

func (m *mockInventoryService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockInventoryService) AdjustQuantity(ctx context.Context, input *service.AdjustQuantityInput) (*model.Inventory, error) {
	return &model.Inventory{}, nil
}

func (m *mockInventoryService) CheckStock(ctx context.Context, input *service.CheckStockInput) (*service.CheckStockResult, error) {
	return &service.CheckStockResult{}, nil
}

type mockSupplierService struct{}

func (m *mockSupplierService) Create(ctx context.Context, input *service.CreateSupplierInput) (*model.Supplier, error) {
	return &model.Supplier{}, nil
}

func (m *mockSupplierService) GetByID(ctx context.Context, id int64) (*model.Supplier, error) {
	return &model.Supplier{}, nil
}

func (m *mockSupplierService) List(ctx context.Context, page, pageSize int, keyword string) (*service.ListSuppliersResult, error) {
	return &service.ListSuppliersResult{}, nil
}

func (m *mockSupplierService) Update(ctx context.Context, id int64, input *service.UpdateSupplierInput) (*model.Supplier, error) {
	return &model.Supplier{}, nil
}

func (m *mockSupplierService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockCustomerService struct{}

func (m *mockCustomerService) Create(ctx context.Context, input *service.CreateCustomerInput) (*model.Customer, error) {
	return &model.Customer{}, nil
}

func (m *mockCustomerService) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	return &model.Customer{}, nil
}

func (m *mockCustomerService) List(ctx context.Context, page, pageSize int, keyword string) (*service.ListCustomersResult, error) {
	return &service.ListCustomersResult{}, nil
}

func (m *mockCustomerService) Update(ctx context.Context, id int64, input *service.UpdateCustomerInput) (*model.Customer, error) {
	return &model.Customer{}, nil
}

func (m *mockCustomerService) Delete(ctx context.Context, id int64) error {
	return nil
}

type mockInboundOrderService struct{}

func (m *mockInboundOrderService) Create(ctx context.Context, input *service.CreateInboundOrderInput) (*model.InboundOrder, error) {
	return &model.InboundOrder{}, nil
}

func (m *mockInboundOrderService) GetByID(ctx context.Context, id int64) (*model.InboundOrder, error) {
	return &model.InboundOrder{}, nil
}

func (m *mockInboundOrderService) List(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListInboundOrdersResult, error) {
	return &service.ListInboundOrdersResult{}, nil
}

func (m *mockInboundOrderService) Update(ctx context.Context, id int64, input *service.UpdateInboundOrderInput) (*model.InboundOrder, error) {
	return &model.InboundOrder{}, nil
}

func (m *mockInboundOrderService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockInboundOrderService) Confirm(ctx context.Context, id int64) (*model.InboundOrder, error) {
	return &model.InboundOrder{}, nil
}

type mockOutboundOrderService struct{}

func (m *mockOutboundOrderService) Create(ctx context.Context, input *service.CreateOutboundOrderInput) (*model.OutboundOrder, error) {
	return &model.OutboundOrder{}, nil
}

func (m *mockOutboundOrderService) GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error) {
	return &model.OutboundOrder{}, nil
}

func (m *mockOutboundOrderService) List(ctx context.Context, page, pageSize int, warehouseID, status int) (*service.ListOutboundOrdersResult, error) {
	return &service.ListOutboundOrdersResult{}, nil
}

func (m *mockOutboundOrderService) Update(ctx context.Context, id int64, input *service.UpdateOutboundOrderInput) (*model.OutboundOrder, error) {
	return &model.OutboundOrder{}, nil
}

func (m *mockOutboundOrderService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockOutboundOrderService) Confirm(ctx context.Context, id int64) (*model.OutboundOrder, error) {
	return &model.OutboundOrder{}, nil
}

type mockStockTransferService struct{}

func (m *mockStockTransferService) Create(ctx context.Context, input *service.CreateStockTransferInput) (*model.StockTransfer, error) {
	return &model.StockTransfer{}, nil
}

func (m *mockStockTransferService) GetByID(ctx context.Context, id int64) (*model.StockTransfer, error) {
	return &model.StockTransfer{}, nil
}

func (m *mockStockTransferService) List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) (*service.ListStockTransfersResult, error) {
	return &service.ListStockTransfersResult{}, nil
}

func (m *mockStockTransferService) Update(ctx context.Context, id int64, input *service.UpdateStockTransferInput) (*model.StockTransfer, error) {
	return &model.StockTransfer{}, nil
}

func (m *mockStockTransferService) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockStockTransferService) Confirm(ctx context.Context, id int64) (*model.StockTransfer, error) {
	return &model.StockTransfer{}, nil
}

type mockAuditLogService struct{}

func (m *mockAuditLogService) GetByID(ctx context.Context, id int64) (*model.AuditLog, error) {
	return &model.AuditLog{}, nil
}

func (m *mockAuditLogService) List(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error) {
	return &service.AuditLogListResult{}, nil
}
