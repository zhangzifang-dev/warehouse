package router

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"warehouse/internal/handler"
	"warehouse/internal/middleware"
	"warehouse/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth          *handler.AuthHandler
	User          *handler.UserHandler
	Role          *handler.RoleHandler
	Permission    *handler.PermissionHandler
	Warehouse     *handler.WarehouseHandler
	Location      *handler.LocationHandler
	Category      *handler.CategoryHandler
	Product       *handler.ProductHandler
	Inventory     *handler.InventoryHandler
	Supplier      *handler.SupplierHandler
	Customer      *handler.CustomerHandler
	InboundOrder  *handler.InboundOrderHandler
	OutboundOrder *handler.OutboundOrderHandler
	StockTransfer *handler.StockTransferHandler
	AuditLog      *handler.AuditLogHandler
}

func Setup(r *gin.Engine, jwtService *jwt.JWT, handlers *Handlers, staticFS embed.FS) {
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Auth.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.JWTAuth(jwtService))
		{
			protected.GET("/auth/profile", handlers.Auth.GetProfile)
			protected.PUT("/auth/password", handlers.Auth.ChangePassword)

			handler.RegisterUserRoutes(protected, handlers.User)
			handler.RegisterRoleRoutes(protected, handlers.Role)
			handler.RegisterPermissionRoutes(protected, handlers.Permission)
			handler.RegisterWarehouseRoutes(protected, handlers.Warehouse)
			handler.RegisterLocationRoutes(protected, handlers.Location)
			handler.RegisterCategoryRoutes(protected, handlers.Category)
			handler.RegisterProductRoutes(protected, handlers.Product)
			handler.RegisterInventoryRoutes(protected, handlers.Inventory)
			handler.RegisterSupplierRoutes(protected, handlers.Supplier)
			handler.RegisterCustomerRoutes(protected, handlers.Customer)
			handler.RegisterInboundOrderRoutes(protected, handlers.InboundOrder)
			handler.RegisterOutboundOrderRoutes(protected, handlers.OutboundOrder)
			handler.RegisterStockTransferRoutes(protected, handlers.StockTransfer)
			handler.RegisterAuditLogRoutes(protected, handlers.AuditLog)
		}
	}

	setupStaticFiles(r, staticFS)
}

func setupStaticFiles(r *gin.Engine, staticFS embed.FS) {
	distFS, err := fs.Sub(staticFS, "web/dist")
	if err != nil {
		return
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		if path == "/" {
			data, err := fs.ReadFile(distFS, "index.html")
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}

		filePath := strings.TrimPrefix(path, "/")
		if _, err := fs.Stat(distFS, filePath); err == nil {
			c.FileFromFS(filePath, http.FS(distFS))
			return
		}

		data, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
}
