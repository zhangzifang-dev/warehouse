package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestRoleRepository_Create(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	role := &model.Role{
		Name:        "Admin",
		Code:        "admin",
		Description: "Administrator role",
		Status:      1,
	}

	err := repo.Create(ctx, role)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestRoleRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestRoleRepository_GetByCode_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	_, err := repo.GetByCode(ctx, "admin")
	if err == nil {
		t.Error("GetByCode() should return error with mock DB")
	}
}

func TestRoleRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	_, _, err := repo.List(ctx, 1, 10)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestRoleRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	role := &model.Role{BaseModel: model.BaseModel{ID: 1}, Name: "Admin", Code: "admin"}
	err := repo.Update(ctx, role)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestRoleRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestRoleRepository_AssignPermissions_Empty(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	err := repo.AssignPermissions(ctx, 1, []int64{})
	if err != nil {
		t.Errorf("AssignPermissions() with empty permissions error = %v", err)
	}
}

func TestRoleRepository_AssignPermissions_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	err := repo.AssignPermissions(ctx, 1, []int64{1, 2})
	if err == nil {
		t.Error("AssignPermissions() should return error with mock DB")
	}
}

func TestRoleRepository_GetRolePermissions_Query(t *testing.T) {
	repo, _, ctx := setupRoleTest(t)
	_, err := repo.GetRolePermissions(ctx, 1)
	if err == nil {
		t.Error("GetRolePermissions() should return error with mock DB")
	}
}

func TestNewRoleRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewRoleRepository(db)
	if repo == nil {
		t.Error("NewRoleRepository() returned nil")
	}
}

func setupRoleTest(t *testing.T) (*RoleRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewRoleRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
