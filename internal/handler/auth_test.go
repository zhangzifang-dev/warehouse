package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	appjwt "warehouse/internal/pkg/jwt"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	loginFunc        func(ctx context.Context, username, password string) (string, *model.User, error)
	getProfileFunc   func(ctx context.Context, userID int64) (*model.User, error)
	changePasswordFunc func(ctx context.Context, userID int64, oldPassword, newPassword string) error
}

func (m *mockAuthService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, username, password)
	}
	return "", nil, errors.New("not implemented")
}

func (m *mockAuthService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	if m.changePasswordFunc != nil {
		return m.changePasswordFunc(ctx, userID, oldPassword, newPassword)
	}
	return errors.New("not implemented")
}

func (m *mockAuthService) UpdateTheme(ctx context.Context, userID int64, theme string) error {
	return nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFunc: func(ctx context.Context, username, password string) (string, *model.User, error) {
			return "test-token", &model.User{
				BaseModel: model.BaseModel{ID: 1},
				Username:  "testuser",
			}, nil
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	body := map[string]string{"username": "testuser", "password": "password"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	if data["token"] != "test-token" {
		t.Errorf("expected token 'test-token', got '%v'", data["token"])
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFunc: func(ctx context.Context, username, password string) (string, *model.User, error) {
			return "", nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	body := map[string]string{"username": "nonexistent", "password": "password"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_InvalidPassword(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFunc: func(ctx context.Context, username, password string) (string, *model.User, error) {
			return "", nil, apperrors.NewAppError(apperrors.CodeInvalidPassword, "invalid password")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	body := map[string]string{"username": "testuser", "password": "wrongpass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_DisabledUser(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFunc: func(ctx context.Context, username, password string) (string, *model.User, error) {
			return "", nil, apperrors.NewAppError(apperrors.CodeForbidden, "user is disabled")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	body := map[string]string{"username": "testuser", "password": "password"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthHandler_Login_BindingError(t *testing.T) {
	mockSvc := &mockAuthService{}
	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_GetProfile_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		getProfileFunc: func(ctx context.Context, userID int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: userID},
				Username:  "testuser",
				Status:    model.UserStatusActive,
			}, nil
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/auth/profile", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	}, handler.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in response")
	}
	if data["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got '%v'", data["username"])
	}
}

func TestAuthHandler_GetProfile_UserNotFound(t *testing.T) {
	mockSvc := &mockAuthService{
		getProfileFunc: func(ctx context.Context, userID int64) (*model.User, error) {
			return nil, apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.GET("/auth/profile", func(c *gin.Context) {
		c.Set("user_id", int64(999))
		c.Next()
	}, handler.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_ChangePassword_Success(t *testing.T) {
	passwordChanged := false
	mockSvc := &mockAuthService{
		changePasswordFunc: func(ctx context.Context, userID int64, oldPassword, newPassword string) error {
			passwordChanged = true
			return nil
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/auth/password", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	}, handler.ChangePassword)

	body := map[string]string{"old_password": "oldpass", "new_password": "newpass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !passwordChanged {
		t.Error("expected password to be changed")
	}
}

func TestAuthHandler_ChangePassword_UserNotFound(t *testing.T) {
	mockSvc := &mockAuthService{
		changePasswordFunc: func(ctx context.Context, userID int64, oldPassword, newPassword string) error {
			return apperrors.NewAppError(apperrors.CodeUserNotFound, "user not found")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/auth/password", func(c *gin.Context) {
		c.Set("user_id", int64(999))
		c.Next()
	}, handler.ChangePassword)

	body := map[string]string{"old_password": "oldpass", "new_password": "newpass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_ChangePassword_InvalidOldPassword(t *testing.T) {
	mockSvc := &mockAuthService{
		changePasswordFunc: func(ctx context.Context, userID int64, oldPassword, newPassword string) error {
			return apperrors.NewAppError(apperrors.CodeInvalidPassword, "invalid old password")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/auth/password", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	}, handler.ChangePassword)

	body := map[string]string{"old_password": "wrongpass", "new_password": "newpass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_ChangePassword_DisabledUser(t *testing.T) {
	mockSvc := &mockAuthService{
		changePasswordFunc: func(ctx context.Context, userID int64, oldPassword, newPassword string) error {
			return apperrors.NewAppError(apperrors.CodeForbidden, "user is disabled")
		},
	}

	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/auth/password", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	}, handler.ChangePassword)

	body := map[string]string{"old_password": "oldpass", "new_password": "newpass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestAuthHandler_ChangePassword_BindingError(t *testing.T) {
	mockSvc := &mockAuthService{}
	handler := NewAuthHandler(mockSvc)
	router := setupTestRouter()
	router.PUT("/auth/password", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	}, handler.ChangePassword)

	req := httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Integration_LoginAndGetProfile(t *testing.T) {
	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	
	mockUserRepo := &mockUserRepositoryForIntegration{
		users: map[int64]*model.User{
			1: {
				BaseModel: model.BaseModel{ID: 1},
				Username:  "testuser",
				Password:  "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe",
				Status:    model.UserStatusActive,
			},
		},
	}
	
	authService := service.NewAuthService(mockUserRepo, jwtSvc)
	handler := NewAuthHandler(authService)
	
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)
	router.GET("/auth/profile", handler.GetProfile)
	
	loginBody := map[string]string{"username": "testuser", "password": "password"}
	jsonBody, _ := json.Marshal(loginBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Fatalf("login failed with status %d", w.Code)
	}
	
	var loginResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	
	data := loginResp["data"].(map[string]interface{})
	token := data["token"].(string)
	
	claims, err := jwtSvc.ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	
	_ = claims
}

type mockUserRepositoryForIntegration struct {
	users map[int64]*model.User
}

func (m *mockUserRepositoryForIntegration) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepositoryForIntegration) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepositoryForIntegration) Update(ctx context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepositoryForIntegration) UpdateTheme(ctx context.Context, userID int64, theme string) error {
	return nil
}
