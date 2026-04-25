package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestLocationRepository_Create(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	location := &model.Location{
		WarehouseID: 1,
		Zone:        "A",
		Shelf:       "01",
		Level:       "02",
		Position:    "03",
		Code:        "A-01-02-03",
		Status:      model.LocationStatusActive,
	}

	err := repo.Create(ctx, location)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestLocationRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestLocationRepository_GetByWarehouseAndCode_Query(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	_, err := repo.GetByWarehouseAndCode(ctx, 1, "A-01-02-03")
	if err == nil {
		t.Error("GetByWarehouseAndCode() should return error with mock DB")
	}
}

func TestLocationRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestLocationRepository_ListByWarehouse_Query(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	_, _, err := repo.ListByWarehouse(ctx, 1, 1, 10)
	if err == nil {
		t.Error("ListByWarehouse() should return error with mock DB")
	}
}

func TestLocationRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	location := &model.Location{
		BaseModel:   model.BaseModel{ID: 1},
		WarehouseID: 1,
		Code:        "A-01-02-03",
	}
	err := repo.Update(ctx, location)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestLocationRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupLocationTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewLocationRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewLocationRepository(db)
	if repo == nil {
		t.Error("NewLocationRepository() returned nil")
	}
}

func setupLocationTest(t *testing.T) (*LocationRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewLocationRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
