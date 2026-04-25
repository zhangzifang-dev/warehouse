package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestPermissionRepository_Create(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	perm := &model.Permission{
		Name:        "Create User",
		Code:        "user:create",
		Resource:    "user",
		Action:      "create",
		Description: "Permission to create users",
	}

	err := repo.Create(ctx, perm)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestPermissionRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestPermissionRepository_GetByCode_Query(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	_, err := repo.GetByCode(ctx, "user:create")
	if err == nil {
		t.Error("GetByCode() should return error with mock DB")
	}
}

func TestPermissionRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	_, _, err := repo.List(ctx, 1, 10)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestPermissionRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	perm := &model.Permission{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Updated Permission",
		Code:        "user:update",
		Resource:    "user",
		Action:      "update",
		Description: "Updated description",
	}
	err := repo.Update(ctx, perm)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestPermissionRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestPermissionRepository_BatchCreate_Empty(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	err := repo.BatchCreate(ctx, []*model.Permission{})
	if err != nil {
		t.Errorf("BatchCreate() with empty permissions error = %v", err)
	}
}

func TestPermissionRepository_BatchCreate_Query(t *testing.T) {
	repo, _, ctx := setupPermissionTest(t)
	perms := []*model.Permission{
		{Name: "Read User", Code: "user:read", Resource: "user", Action: "read"},
		{Name: "Update User", Code: "user:update", Resource: "user", Action: "update"},
	}
	err := repo.BatchCreate(ctx, perms)
	if err == nil {
		t.Error("BatchCreate() should return error with mock DB")
	}
}

func TestNewPermissionRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewPermissionRepository(db)
	if repo == nil {
		t.Error("NewPermissionRepository() returned nil")
	}
}

func setupPermissionTest(t *testing.T) (*PermissionRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewPermissionRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
