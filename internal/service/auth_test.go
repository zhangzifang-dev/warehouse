package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"warehouse/internal/model"
	appjwt "warehouse/internal/pkg/jwt"
)

type mockUserRepository struct {
	getByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
	getByIDFunc       func(ctx context.Context, id int64) (*model.User, error)
	updateFunc        func(ctx context.Context, user *model.User) error
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.getByUsernameFunc != nil {
		return m.getByUsernameFunc(ctx, username)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return errors.New("not implemented")
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: 1},
				Username:  "testuser",
				Password:  "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe",
				Status:    model.UserStatusActive,
			}, nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	token, user, err := svc.Login(context.Background(), "testuser", "password")

	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Error("expected token, got empty string")
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	_, _, err := svc.Login(context.Background(), "nonexistent", "password")

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: 1},
				Username:  "testuser",
				Password:  "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe",
				Status:    model.UserStatusActive,
			}, nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	_, _, err := svc.Login(context.Background(), "testuser", "wrongpassword")

	if err == nil {
		t.Error("expected error for invalid password, got nil")
	}
}

func TestAuthService_Login_DisabledUser(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByUsernameFunc: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: 1},
				Username:  "testuser",
				Password:  "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe",
				Status:    model.UserStatusDisabled,
			}, nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	_, _, err := svc.Login(context.Background(), "testuser", "password")

	if err == nil {
		t.Error("expected error for disabled user, got nil")
	}
}

func TestAuthService_GetProfile_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Nickname:  "Test User",
				Email:     "test@example.com",
				Phone:     "1234567890",
				Status:    model.UserStatusActive,
			}, nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	user, err := svc.GetProfile(context.Background(), 1)

	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}
}

func TestAuthService_GetProfile_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	_, err := svc.GetProfile(context.Background(), 999)

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestAuthService_ChangePassword_Success(t *testing.T) {
	oldHash := "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe"
	updatedUser := &model.User{}

	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Password:  oldHash,
				Status:    model.UserStatusActive,
			}, nil
		},
		updateFunc: func(ctx context.Context, user *model.User) error {
			updatedUser = user
			return nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	err := svc.ChangePassword(context.Background(), 1, "password", "newpassword")

	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}
	if updatedUser.Password == oldHash {
		t.Error("expected password to be updated, but it was not changed")
	}
}

func TestAuthService_ChangePassword_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	err := svc.ChangePassword(context.Background(), 999, "oldpass", "newpass")

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}

func TestAuthService_ChangePassword_InvalidOldPassword(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Password:  "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe",
				Status:    model.UserStatusActive,
			}, nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	err := svc.ChangePassword(context.Background(), 1, "wrongpassword", "newpassword")

	if err == nil {
		t.Error("expected error for invalid old password, got nil")
	}
}

func TestAuthService_ChangePassword_DisabledUser(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{
				BaseModel: model.BaseModel{ID: id},
				Username:  "testuser",
				Password:  "$2a$10$L3paV72KbqXLrVeEaCBNVODAIR661qEjPgB5Em8815WSd19uljYfe",
				Status:    model.UserStatusDisabled,
			}, nil
		},
	}

	jwtSvc := appjwt.NewJWT("test-secret-key", time.Hour)
	svc := NewAuthService(mockRepo, jwtSvc)

	err := svc.ChangePassword(context.Background(), 1, "password", "newpassword")

	if err == nil {
		t.Error("expected error for disabled user, got nil")
	}
}
