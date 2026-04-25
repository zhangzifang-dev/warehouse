package middleware

import (
	"context"
	"log"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	UserIDKey       contextKey = "userID"
	PermissionsKey  contextKey = "permissions"
)

type RBACUserRepository interface {
	GetUserPermissions(ctx context.Context, userID int64) ([]model.Permission, error)
}

func RBACAuth(userRepo RBACUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(string(UserIDKey))
		if !exists {
			response.Error(c, 401, "unauthorized")
			c.Abort()
			return
		}

		uid, ok := userID.(int64)
		if !ok {
			response.Error(c, 401, "invalid user id")
			c.Abort()
			return
		}

		permissions, err := userRepo.GetUserPermissions(c.Request.Context(), uid)
		if err != nil {
			log.Printf("failed to get user permissions: %v", err)
			response.Error(c, 500, "failed to load permissions")
			c.Abort()
			return
		}

		codes := make([]string, len(permissions))
		for i, perm := range permissions {
			codes[i] = perm.Code
		}

		c.Set(string(PermissionsKey), codes)
		c.Next()
	}
}

func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, exists := c.Get(string(PermissionsKey))
		if !exists {
			response.Error(c, 401, "unauthorized")
			c.Abort()
			return
		}

		codes, ok := perms.([]string)
		if !ok {
			response.Error(c, 500, "invalid permissions")
			c.Abort()
			return
		}

		for _, permCode := range codes {
			if permCode == code {
				c.Next()
				return
			}
		}

		response.Error(c, 403, "insufficient permissions")
		c.Abort()
	}
}
