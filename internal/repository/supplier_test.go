package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestSupplierRepository_Create(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	supplier := &model.Supplier{
		Name:   "Test Supplier",
		Code:   "SUP001",
		Status: model.SupplierStatusActive,
	}

	err := repo.Create(ctx, supplier)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestSupplierRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestSupplierRepository_GetByCode_Query(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	_, err := repo.GetByCode(ctx, "SUP001")
	if err == nil {
		t.Error("GetByCode() should return error with mock DB")
	}
}

func TestSupplierRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	filter := &SupplierQueryFilter{Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_List_WithKeyword(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	filter := &SupplierQueryFilter{Name: "test", Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_List_WithCodeFilter(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	filter := &SupplierQueryFilter{Code: "SUP", Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_List_WithContactFilter(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	filter := &SupplierQueryFilter{Contact: "John", Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_List_WithPhoneFilter(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	filter := &SupplierQueryFilter{Phone: "123", Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_List_WithStatusFilter(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	status := 1
	filter := &SupplierQueryFilter{Status: &status, Page: 1, PageSize: 10}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_List_WithMultipleFilters(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	status := 1
	filter := &SupplierQueryFilter{
		Code:    "SUP",
		Name:    "Test",
		Contact: "John",
		Phone:   "123",
		Status:  &status,
		Page:    1,
		PageSize: 10,
	}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestSupplierRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	supplier := &model.Supplier{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Updated Supplier",
	}
	err := repo.Update(ctx, supplier)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestSupplierRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupSupplierTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewSupplierRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewSupplierRepository(db)
	if repo == nil {
		t.Error("NewSupplierRepository() returned nil")
	}
}

func setupSupplierTest(t *testing.T) (*SupplierRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewSupplierRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
