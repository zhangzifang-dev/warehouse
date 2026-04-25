package middleware

import (
	"strings"

	"warehouse/internal/pkg/errors"
	"warehouse/internal/pkg/jwt"
	"warehouse/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
)

func JWTAuth(jwtService *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.CodeUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, errors.CodeUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			response.Error(c, errors.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)

		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0
	}
	return userID.(int64)
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return ""
	}
	return username.(string)
}
