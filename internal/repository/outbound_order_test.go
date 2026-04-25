package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestOutboundOrderRepository_Create(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	order := &model.OutboundOrder{
		OrderNo:       "SO-2024-001",
		WarehouseID:   1,
		TotalQuantity: 100,
		Status:        0,
	}
	err := repo.Create(ctx, order)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_GetByID(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_GetByOrderNo(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	_, err := repo.GetByOrderNo(ctx, "SO-2024-001")
	if err == nil {
		t.Error("GetByOrderNo() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_List(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0, 0)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_List_WithWarehouseID(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	_, _, err := repo.List(ctx, 1, 10, 1, 0)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_List_WithStatus(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0, 1)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_Update(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	order := &model.OutboundOrder{
		BaseModel:     model.BaseModel{ID: 1},
		OrderNo:       "SO-2024-001",
		WarehouseID:   1,
		TotalQuantity: 200,
		Status:        1,
	}
	err := repo.Update(ctx, order)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestOutboundOrderRepository_Delete(t *testing.T) {
	repo, _, ctx := setupOutboundOrderTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewOutboundOrderRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewOutboundOrderRepository(db)
	if repo == nil {
		t.Error("NewOutboundOrderRepository() returned nil")
	}
}

func setupOutboundOrderTest(t *testing.T) (*OutboundOrderRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewOutboundOrderRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
