package main

import (
	"fmt"
	"log"
	"time"

	embedfs "warehouse"
	"warehouse/internal/config"
	"warehouse/internal/handler"
	"warehouse/internal/pkg/jwt"
	"warehouse/internal/repository"
	"warehouse/internal/router"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	inboundOrderRepo := repository.NewInboundOrderRepository(db)
	inboundItemRepo := repository.NewInboundItemRepository(db)
	outboundOrderRepo := repository.NewOutboundOrderRepository(db)
	outboundItemRepo := repository.NewOutboundItemRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	expireDuration, err := time.ParseDuration(cfg.JWT.Expire)
	if err != nil {
		expireDuration = 24 * time.Hour
	}
	jwtService := jwt.NewJWT(cfg.JWT.Secret, expireDuration)

	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	locationService := service.NewLocationService(locationRepo, warehouseRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	supplierService := service.NewSupplierService(supplierRepo)
	customerService := service.NewCustomerService(customerRepo)
	inboundOrderService := service.NewInboundOrderService(inboundOrderRepo, inboundItemRepo, inventoryService)
	outboundOrderService := service.NewOutboundOrderService(outboundOrderRepo, outboundItemRepo, inventoryService)
	stockTransferService := service.NewStockTransferService(nil, nil, inventoryService)
	auditLogService := service.NewAuditLogService(auditLogRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)
	locationHandler := handler.NewLocationHandler(locationService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	customerHandler := handler.NewCustomerHandler(customerService)
	inboundOrderHandler := handler.NewInboundOrderHandler(inboundOrderService)
	outboundOrderHandler := handler.NewOutboundOrderHandler(outboundOrderService)
	stockTransferHandler := handler.NewStockTransferHandler(stockTransferService)
	auditLogHandler := handler.NewAuditLogHandler(auditLogService)

	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()

	router.Setup(engine, jwtService, &router.Handlers{
		Auth:          authHandler,
		User:          userHandler,
		Role:          roleHandler,
		Permission:    permissionHandler,
		Warehouse:     warehouseHandler,
		Location:      locationHandler,
		Category:      categoryHandler,
		Product:       productHandler,
		Inventory:     inventoryHandler,
		Supplier:      supplierHandler,
		Customer:      customerHandler,
		InboundOrder:  inboundOrderHandler,
		OutboundOrder: outboundOrderHandler,
		StockTransfer: stockTransferHandler,
		AuditLog:      auditLogHandler,
	}, embedfs.Static)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
