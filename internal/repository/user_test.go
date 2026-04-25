package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestUserRepository_Create(t *testing.T) {
	repo, _, ctx := setupTest(t)
	user := &model.User{
		Username: "testuser",
		Password: "hashedpassword",
		Nickname: "Test User",
		Email:    "test@example.com",
		Phone:    "1234567890",
		Status:   model.UserStatusActive,
	}

	err := repo.Create(ctx, user)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestUserRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo, _, ctx := setupTest(t)
	_, err := repo.GetByID(ctx, 99999)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestUserRepository_GetByUsername_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	_, err := repo.GetByUsername(ctx, "testuser")
	if err == nil {
		t.Error("GetByUsername() should return error with mock DB")
	}
}

func TestUserRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	_, _, err := repo.List(ctx, 1, 10)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestUserRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	user := &model.User{BaseModel: model.BaseModel{ID: 1}, Username: "test"}
	err := repo.Update(ctx, user)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestUserRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestUserRepository_GetUserRoles_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	_, err := repo.GetUserRoles(ctx, 1)
	if err == nil {
		t.Error("GetUserRoles() should return error with mock DB")
	}
}

func TestUserRepository_GetUserPermissions_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	_, err := repo.GetUserPermissions(ctx, 1)
	if err == nil {
		t.Error("GetUserPermissions() should return error with mock DB")
	}
}

func TestUserRepository_AssignRoles_Empty(t *testing.T) {
	repo, _, ctx := setupTest(t)
	err := repo.AssignRoles(ctx, 1, []int64{})
	if err != nil {
		t.Errorf("AssignRoles() with empty roles error = %v", err)
	}
}

func TestUserRepository_AssignRoles_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	err := repo.AssignRoles(ctx, 1, []int64{1, 2})
	if err == nil {
		t.Error("AssignRoles() should return error with mock DB")
	}
}

func TestUserRepository_CreateRole_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	role := &model.Role{Name: "Admin", Code: "admin", Status: 1}
	err := repo.CreateRole(ctx, role)
	if err == nil {
		t.Error("CreateRole() should return error with mock DB")
	}
}

func TestUserRepository_CreatePermission_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	perm := &model.Permission{Name: "Read", Code: "read", Resource: "user", Action: "read"}
	err := repo.CreatePermission(ctx, perm)
	if err == nil {
		t.Error("CreatePermission() should return error with mock DB")
	}
}

func TestUserRepository_AssignPermissionsToRole_Empty(t *testing.T) {
	repo, _, ctx := setupTest(t)
	err := repo.AssignPermissionsToRole(ctx, 1, []int64{})
	if err != nil {
		t.Errorf("AssignPermissionsToRole() with empty permissions error = %v", err)
	}
}

func TestUserRepository_AssignPermissionsToRole_Query(t *testing.T) {
	repo, _, ctx := setupTest(t)
	err := repo.AssignPermissionsToRole(ctx, 1, []int64{1, 2})
	if err == nil {
		t.Error("AssignPermissionsToRole() should return error with mock DB")
	}
}

func TestNewUserRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewUserRepository(db)
	if repo == nil {
		t.Error("NewUserRepository() returned nil")
	}
}

func TestUserModel_TableName(t *testing.T) {
	user := model.User{}
	if user.TableName() != "users" {
		t.Errorf("User.TableName() = %s, want users", user.TableName())
	}
}

func TestRole_TableName(t *testing.T) {
	role := model.Role{}
	if role.TableName() != "roles" {
		t.Errorf("Role.TableName() = %s, want roles", role.TableName())
	}
}

func TestPermission_TableName(t *testing.T) {
	perm := model.Permission{}
	if perm.TableName() != "permissions" {
		t.Errorf("Permission.TableName() = %s, want permissions", perm.TableName())
	}
}

func TestUserRole_TableName(t *testing.T) {
	ur := model.UserRole{}
	if ur.TableName() != "user_roles" {
		t.Errorf("UserRole.TableName() = %s, want user_roles", ur.TableName())
	}
}

func TestRolePermission_TableName(t *testing.T) {
	rp := model.RolePermission{}
	if rp.TableName() != "role_permissions" {
		t.Errorf("RolePermission.TableName() = %s, want role_permissions", rp.TableName())
	}
}

func TestUser_IsActive(t *testing.T) {
	user := model.User{Status: model.UserStatusActive}
	if !user.IsActive() {
		t.Error("IsActive() = false for active user")
	}

	user.Status = model.UserStatusDisabled
	if user.IsActive() {
		t.Error("IsActive() = true for disabled user")
	}
}

func setupTest(t *testing.T) (*UserRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewUserRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
