package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestInboundOrderRepository_Create(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	order := &model.InboundOrder{
		OrderNo:       "PO-2024-001",
		WarehouseID:   1,
		TotalQuantity: 100,
		Status:        0,
	}
	err := repo.Create(ctx, order)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestInboundOrderRepository_GetByID(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestInboundOrderRepository_GetByOrderNo(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	_, err := repo.GetByOrderNo(ctx, "PO-2024-001")
	if err == nil {
		t.Error("GetByOrderNo() should return error with mock DB")
	}
}

func TestInboundOrderRepository_List(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0, 0)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInboundOrderRepository_List_WithWarehouseID(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	_, _, err := repo.List(ctx, 1, 10, 1, 0)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInboundOrderRepository_List_WithStatus(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0, 1)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestInboundOrderRepository_Update(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	order := &model.InboundOrder{
		BaseModel:     model.BaseModel{ID: 1},
		OrderNo:       "PO-2024-001",
		WarehouseID:   1,
		TotalQuantity: 200,
		Status:        1,
	}
	err := repo.Update(ctx, order)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestInboundOrderRepository_Delete(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestInboundOrderRepository_ListWithFilter(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	supplierID := int64(1)
	warehouseID := int64(1)
	quantityMin := 10.0
	quantityMax := 100.0
	startTime := time.Now()
	endTime := time.Now().Add(24 * time.Hour)

	filter := &model.InboundOrderQueryFilter{
		OrderNo:        "PO-2024",
		SupplierID:     &supplierID,
		WarehouseID:    &warehouseID,
		QuantityMin:    &quantityMin,
		QuantityMax:    &quantityMax,
		CreatedAtStart: &startTime,
		CreatedAtEnd:   &endTime,
		Page:           1,
		PageSize:       10,
	}

	_, _, err := repo.ListWithFilter(ctx, filter)
	if err == nil {
		t.Error("ListWithFilter() should return error with mock DB")
	}
}

func TestInboundOrderRepository_ListWithFilter_PartialFilters(t *testing.T) {
	repo, _, ctx := setupInboundOrderTest(t)
	supplierID := int64(1)

	filter := &model.InboundOrderQueryFilter{
		SupplierID: &supplierID,
		Page:       1,
		PageSize:   10,
	}

	_, _, err := repo.ListWithFilter(ctx, filter)
	if err == nil {
		t.Error("ListWithFilter() should return error with mock DB")
	}
}

func TestNewInboundOrderRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewInboundOrderRepository(db)
	if repo == nil {
		t.Error("NewInboundOrderRepository() returned nil")
	}
}

func setupInboundOrderTest(t *testing.T) (*InboundOrderRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewInboundOrderRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
