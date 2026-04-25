package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"warehouse/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	permissions []model.Permission
	err         error
}

func (m *mockUserRepository) GetUserPermissions(ctx context.Context, userID int64) ([]model.Permission, error) {
	return m.permissions, m.err
}

func TestRBACAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 401 when user id not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		handler := RBACAuth(nil)
		handler(c)

		assert.True(t, c.IsAborted())
	})

	t.Run("returns 401 when user id is not int64", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Set(string(UserIDKey), "not-int64")

		handler := RBACAuth(nil)
		handler(c)

		assert.True(t, c.IsAborted())
	})

	t.Run("loads permissions into context", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			permissions: []model.Permission{
				{Code: "user:read"},
				{Code: "user:write"},
			},
		}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Set(string(UserIDKey), int64(1))

		handler := RBACAuth(mockRepo)
		handler(c)

		assert.False(t, c.IsAborted())
		perms, exists := c.Get(string(PermissionsKey))
		assert.True(t, exists)
		codes := perms.([]string)
		assert.Equal(t, []string{"user:read", "user:write"}, codes)
	})
}

func TestRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 401 when permissions not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		handler := RequirePermission("user:read")
		handler(c)

		assert.True(t, c.IsAborted())
	})

	t.Run("returns 403 when permission not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Set(string(PermissionsKey), []string{"user:read"})

		handler := RequirePermission("user:delete")
		handler(c)

		assert.True(t, c.IsAborted())
	})

	t.Run("allows request when permission found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Set(string(PermissionsKey), []string{"user:read", "user:write"})

		handler := RequirePermission("user:write")
		handler(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("allows request with exact permission match", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Set(string(PermissionsKey), []string{"admin:all"})

		handler := RequirePermission("admin:all")
		handler(c)

		assert.False(t, c.IsAborted())
	})
}

func TestRBACMiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("full chain allows access with correct permission", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			permissions: []model.Permission{
				{Code: "product:create"},
			},
		}

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(string(UserIDKey), int64(1))
			c.Next()
		})
		r.Use(RBACAuth(mockRepo))
		r.Use(RequirePermission("product:create"))
		r.GET("/products", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("full chain denies access without permission", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			permissions: []model.Permission{
				{Code: "product:read"},
			},
		}

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(string(UserIDKey), int64(1))
			c.Next()
		})
		r.Use(RBACAuth(mockRepo))
		r.Use(RequirePermission("product:delete"))
		r.DELETE("/products/1", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
