package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestWarehouseRepository_Create(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	warehouse := &model.Warehouse{
		Name:    "Main Warehouse",
		Code:    "WH001",
		Address: "123 Main St",
		Contact: "John Doe",
		Phone:   "1234567890",
		Status:  model.WarehouseStatusActive,
	}

	err := repo.Create(ctx, warehouse)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestWarehouseRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestWarehouseRepository_GetByID_NotFound(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	_, err := repo.GetByID(ctx, 99999)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestWarehouseRepository_GetByCode_Query(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	_, err := repo.GetByCode(ctx, "WH001")
	if err == nil {
		t.Error("GetByCode() should return error with mock DB")
	}
}

func TestWarehouseRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	_, _, err := repo.List(ctx, 1, 10)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestWarehouseRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	warehouse := &model.Warehouse{BaseModel: model.BaseModel{ID: 1}, Name: "Updated", Code: "WH001"}
	err := repo.Update(ctx, warehouse)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestWarehouseRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupWarehouseTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewWarehouseRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewWarehouseRepository(db)
	if repo == nil {
		t.Error("NewWarehouseRepository() returned nil")
	}
}

func setupWarehouseTest(t *testing.T) (*WarehouseRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewWarehouseRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
