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

func TestStockTransferRepository_ListWithFilter(t *testing.T) {
	repo, _, ctx := setupStockTransferTest(t)
	sourceWarehouseID := int64(1)
	targetWarehouseID := int64(2)
	startTime := time.Now()
	endTime := time.Now().Add(24 * time.Hour)

	filter := &model.StockTransferQueryFilter{
		OrderNo:           "ST-2024",
		SourceWarehouseID: &sourceWarehouseID,
		TargetWarehouseID: &targetWarehouseID,
		CreatedAtStart:    &startTime,
		CreatedAtEnd:      &endTime,
		Page:              1,
		PageSize:          10,
	}

	_, _, err := repo.ListWithFilter(ctx, filter)
	if err == nil {
		t.Error("ListWithFilter() should return error with mock DB")
	}
}

func TestStockTransferRepository_ListWithFilter_PartialFilters(t *testing.T) {
	repo, _, ctx := setupStockTransferTest(t)
	sourceWarehouseID := int64(1)

	filter := &model.StockTransferQueryFilter{
		SourceWarehouseID: &sourceWarehouseID,
		Page:              1,
		PageSize:          10,
	}

	_, _, err := repo.ListWithFilter(ctx, filter)
	if err == nil {
		t.Error("ListWithFilter() should return error with mock DB")
	}
}

func setupStockTransferTest(t *testing.T) (*StockTransferRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewStockTransferRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
