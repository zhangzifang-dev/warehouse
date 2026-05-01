package service

import (
	"context"
)

type contextKey string

const (
	clientIPKey contextKey = "client_ip"
	userIDKey   contextKey = "user_id"
)

func GetClientIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(clientIPKey).(string); ok {
		return ip
	}
	return ""
}

func SetClientIPToContext(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

func GetUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(userIDKey).(int64); ok {
		return userID
	}
	return 0
}

func SetUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
