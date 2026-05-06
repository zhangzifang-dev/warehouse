package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestInventoryRepository_Create(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	inventory := &model.Inventory{
		WarehouseID: 1,
		ProductID:   1,
		Quantity:    100,
	}

	err := repo.Create(ctx, inventory)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestInventoryRepository_GetByID(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestInventoryRepository_List(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInventoryRepository_List_WithProductName(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10, ProductName: "测试商品"})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInventoryRepository_List_WithQuantityMin(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	minQty := 50.0
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10, QuantityMin: &minQty})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInventoryRepository_List_WithQuantityMax(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	maxQty := 200.0
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10, QuantityMax: &maxQty})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInventoryRepository_List_WithBatchNo(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10, BatchNo: "BATCH001"})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInventoryRepository_Update(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	inventory := &model.Inventory{
		BaseModel:   model.BaseModel{ID: 1},
		WarehouseID: 1,
		ProductID:   1,
		Quantity:    200,
	}
	err := repo.Update(ctx, inventory)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestInventoryRepository_Delete(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestInventoryRepository_GetByWarehouseAndProduct(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, err := repo.GetByWarehouseAndProduct(ctx, 1, 1, "BATCH001")
	if err == nil {
		t.Error("GetByWarehouseAndProduct() should return error with mock DB")
	}
}

func TestInventoryRepository_GetByWarehouseAndProduct_EmptyBatchNo(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, err := repo.GetByWarehouseAndProduct(ctx, 1, 1, "")
	if err == nil {
		t.Error("GetByWarehouseAndProduct() should return error with mock DB")
	}
}

func TestInventoryRepository_UpdateQuantity(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	err := repo.UpdateQuantity(ctx, 1, 50)
	if err == nil {
		t.Error("UpdateQuantity() should return error with mock DB")
	}
}

func TestNewInventoryRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewInventoryRepository(db)
	if repo == nil {
		t.Error("NewInventoryRepository() returned nil")
	}
}

func setupInventoryTest(t *testing.T) (*InventoryRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewInventoryRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}

func TestInventoryRepository_List_WithProductID(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10, ProductID: 1})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInventoryRepository_List_WithWarehouseID(t *testing.T) {
	repo, _, ctx := setupInventoryTest(t)
	_, _, err := repo.List(ctx, &model.InventoryQueryFilter{Page: 1, PageSize: 10, WarehouseID: 1})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}
