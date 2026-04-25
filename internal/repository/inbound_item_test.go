package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestInboundItemRepository_Create(t *testing.T) {
	repo, _, ctx := setupInboundItemTest(t)
	item := &model.InboundItem{
		OrderID:  1,
		ProductID: 1,
		Quantity:  100,
	}
	err := repo.Create(ctx, item)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestInboundItemRepository_ListByOrderID(t *testing.T) {
	repo, _, ctx := setupInboundItemTest(t)
	_, err := repo.ListByOrderID(ctx, 1)
	if err == nil {
		t.Error("ListByOrderID() should return error with mock DB")
	}
}

func TestInboundItemRepository_Update(t *testing.T) {
	repo, _, ctx := setupInboundItemTest(t)
	item := &model.InboundItem{
		BaseModel:  model.BaseModel{ID: 1},
		OrderID:    1,
		ProductID:  1,
		Quantity:   200,
	}
	err := repo.Update(ctx, item)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestInboundItemRepository_Delete(t *testing.T) {
	repo, _, ctx := setupInboundItemTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewInboundItemRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewInboundItemRepository(db)
	if repo == nil {
		t.Error("NewInboundItemRepository() returned nil")
	}
}

func setupInboundItemTest(t *testing.T) (*InboundItemRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewInboundItemRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
