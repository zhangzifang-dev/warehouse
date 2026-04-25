package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"warehouse/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	jwtService := jwt.NewJWT("test-secret", time.Hour)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler := JWTAuth(jwtService)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}
}

func TestJWTAuth_InvalidFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"NoBearer", "token123"},
		{"WrongPrefix", "Basic token123"},
		{"EmptyBearer", "Bearer "},
		{"NoSpace", "Bearertoken"},
	}

	jwtService := jwt.NewJWT("test-secret", time.Hour)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			c.Request.Header.Set("Authorization", tt.header)

			handler := JWTAuth(jwtService)
			handler(c)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
			}
			if !c.IsAborted() {
				t.Error("Request should be aborted")
			}
		})
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	jwtService := jwt.NewJWT("test-secret", time.Hour)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	handler := JWTAuth(jwtService)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if !c.IsAborted() {
		t.Error("Request should be aborted")
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	jwtService := jwt.NewJWT("test-secret", time.Hour)

	userID := int64(123)
	username := "testuser"
	token, err := jwtService.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	called := false
	handler := JWTAuth(jwtService)
	handler(c)

	if c.IsAborted() {
		t.Error("Request should not be aborted")
	}

	_ = called

	gotUserID := GetUserID(c)
	gotUsername := GetUsername(c)

	if gotUserID != userID {
		t.Errorf("UserID = %d, want %d", gotUserID, userID)
	}
	if gotUsername != username {
		t.Errorf("Username = %s, want %s", gotUsername, username)
	}
}

func TestJWTAuth_CaseInsensitiveBearer(t *testing.T) {
	jwtService := jwt.NewJWT("test-secret", time.Hour)

	userID := int64(456)
	username := "testuser2"
	token, err := jwtService.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	tests := []string{
		"Bearer " + token,
		"bearer " + token,
		"BEARER " + token,
		"BeArEr " + token,
	}

	for i, authHeader := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			c.Request.Header.Set("Authorization", authHeader)

			handler := JWTAuth(jwtService)
			handler(c)

			if c.IsAborted() {
				t.Error("Request should not be aborted")
			}

			gotUserID := GetUserID(c)
			if gotUserID != userID {
				t.Errorf("UserID = %d, want %d", gotUserID, userID)
			}
		})
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	jwtService := jwt.NewJWT("test-secret", -time.Hour)

	token, err := jwtService.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := JWTAuth(jwtService)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if !c.IsAborted() {
		t.Error("Request should be aborted for expired token")
	}
}

func TestJWTAuth_WrongSecret(t *testing.T) {
	jwtService1 := jwt.NewJWT("secret1", time.Hour)
	jwtService2 := jwt.NewJWT("secret2", time.Hour)

	token, err := jwtService1.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := JWTAuth(jwtService2)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if !c.IsAborted() {
		t.Error("Request should be aborted for wrong secret")
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID := GetUserID(c)
	if userID != 0 {
		t.Errorf("UserID = %d, want 0", userID)
	}
}

func TestGetUsername_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	username := GetUsername(c)
	if username != "" {
		t.Errorf("Username = %s, want empty", username)
	}
}
