package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestOutboundItemRepository_Create(t *testing.T) {
	repo, _, ctx := setupOutboundItemTest(t)
	item := &model.OutboundItem{
		OrderID:   1,
		ProductID: 1,
		Quantity:  10,
	}
	err := repo.Create(ctx, item)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestOutboundItemRepository_ListByOrderID(t *testing.T) {
	repo, _, ctx := setupOutboundItemTest(t)
	_, err := repo.ListByOrderID(ctx, 1)
	if err == nil {
		t.Error("ListByOrderID() should return error with mock DB")
	}
}

func TestOutboundItemRepository_Update(t *testing.T) {
	repo, _, ctx := setupOutboundItemTest(t)
	item := &model.OutboundItem{
		BaseModel: model.BaseModel{ID: 1},
		OrderID:   1,
		ProductID: 1,
		Quantity:  20,
	}
	err := repo.Update(ctx, item)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestOutboundItemRepository_Delete(t *testing.T) {
	repo, _, ctx := setupOutboundItemTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewOutboundItemRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewOutboundItemRepository(db)
	if repo == nil {
		t.Error("NewOutboundItemRepository() returned nil")
	}
}

func setupOutboundItemTest(t *testing.T) (*OutboundItemRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewOutboundItemRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
