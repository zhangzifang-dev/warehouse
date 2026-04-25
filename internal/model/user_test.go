package model

import (
	"testing"
	"time"
)

func TestUserFields(t *testing.T) {
	now := time.Now()
	user := User{
		BaseModel: BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: 100,
			UpdatedBy: 100,
		},
		Username: "testuser",
		Password: "hashedpassword",
		Nickname: "Test User",
		Email:    "test@example.com",
		Phone:    "1234567890",
		Status:   1,
	}

	if user.ID != 1 {
		t.Errorf("ID = %d, want 1", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", user.Username)
	}
	if user.Password != "hashedpassword" {
		t.Errorf("Password = %s, want hashedpassword", user.Password)
	}
	if user.Nickname != "Test User" {
		t.Errorf("Nickname = %s, want Test User", user.Nickname)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", user.Email)
	}
	if user.Phone != "1234567890" {
		t.Errorf("Phone = %s, want 1234567890", user.Phone)
	}
	if user.Status != 1 {
		t.Errorf("Status = %d, want 1", user.Status)
	}
}

func TestUserZeroValues(t *testing.T) {
	user := User{}

	if user.ID != 0 {
		t.Errorf("ID = %d, want 0", user.ID)
	}
	if user.Username != "" {
		t.Errorf("Username = %s, want empty", user.Username)
	}
	if user.Password != "" {
		t.Errorf("Password = %s, want empty", user.Password)
	}
	if user.Nickname != "" {
		t.Errorf("Nickname = %s, want empty", user.Nickname)
	}
	if user.Email != "" {
		t.Errorf("Email = %s, want empty", user.Email)
	}
	if user.Phone != "" {
		t.Errorf("Phone = %s, want empty", user.Phone)
	}
	if user.Status != 0 {
		t.Errorf("Status = %d, want 0", user.Status)
	}
}

func TestUserStatus(t *testing.T) {
	user := User{Status: UserStatusActive}
	if user.Status != UserStatusActive {
		t.Errorf("Status = %d, want %d", user.Status, UserStatusActive)
	}

	user.Status = UserStatusDisabled
	if user.Status != UserStatusDisabled {
		t.Errorf("Status = %d, want %d", user.Status, UserStatusDisabled)
	}
}

func TestUserIsActive(t *testing.T) {
	user := User{Status: UserStatusActive}
	if !user.IsActive() {
		t.Error("IsActive() = false, want true for active user")
	}

	user.Status = UserStatusDisabled
	if user.IsActive() {
		t.Error("IsActive() = true, want false for disabled user")
	}
}

func TestRoleFields(t *testing.T) {
	now := time.Now()
	role := Role{
		BaseModel: BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Admin",
		Code:        "admin",
		Description: "Administrator role",
		Status:      1,
	}

	if role.Name != "Admin" {
		t.Errorf("Name = %s, want Admin", role.Name)
	}
	if role.Code != "admin" {
		t.Errorf("Code = %s, want admin", role.Code)
	}
	if role.Description != "Administrator role" {
		t.Errorf("Description = %s, want Administrator role", role.Description)
	}
	if role.Status != 1 {
		t.Errorf("Status = %d, want 1", role.Status)
	}
}

func TestPermissionFields(t *testing.T) {
	now := time.Now()
	perm := Permission{
		BaseModel: BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Create User",
		Code:        "user:create",
		Resource:    "user",
		Action:      "create",
		Description: "Permission to create users",
	}

	if perm.Name != "Create User" {
		t.Errorf("Name = %s, want Create User", perm.Name)
	}
	if perm.Code != "user:create" {
		t.Errorf("Code = %s, want user:create", perm.Code)
	}
	if perm.Resource != "user" {
		t.Errorf("Resource = %s, want user", perm.Resource)
	}
	if perm.Action != "create" {
		t.Errorf("Action = %s, want create", perm.Action)
	}
	if perm.Description != "Permission to create users" {
		t.Errorf("Description = %s, want Permission to create users", perm.Description)
	}
}

func TestUserRoleFields(t *testing.T) {
	ur := UserRole{
		BaseModel: BaseModel{ID: 1},
		UserID:    100,
		RoleID:    200,
	}

	if ur.UserID != 100 {
		t.Errorf("UserID = %d, want 100", ur.UserID)
	}
	if ur.RoleID != 200 {
		t.Errorf("RoleID = %d, want 200", ur.RoleID)
	}
}

func TestRolePermissionFields(t *testing.T) {
	rp := RolePermission{
		BaseModel:     BaseModel{ID: 1},
		RoleID:        100,
		PermissionID:  200,
	}

	if rp.RoleID != 100 {
		t.Errorf("RoleID = %d, want 100", rp.RoleID)
	}
	if rp.PermissionID != 200 {
		t.Errorf("PermissionID = %d, want 200", rp.PermissionID)
	}
}
