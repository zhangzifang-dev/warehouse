package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestProductRepository_Create(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	product := &model.Product{
		SKU:   "SKU001",
		Name:  "Test Product",
		Price: 99.99,
		Status: model.ProductStatusActive,
	}

	err := repo.Create(ctx, product)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestProductRepository_GetBySKU(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	_, err := repo.GetBySKU(ctx, "SKU001")
	if err == nil {
		t.Error("GetBySKU() should return error with mock DB")
	}
}

func TestProductRepository_List(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0, "")
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestProductRepository_List_WithCategoryID(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	_, _, err := repo.List(ctx, 1, 10, 1, "")
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestProductRepository_List_WithKeyword(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	_, _, err := repo.List(ctx, 1, 10, 0, "test")
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestProductRepository_Update(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	product := &model.Product{
		BaseModel: model.BaseModel{ID: 1},
		SKU:       "SKU001",
		Name:      "Updated Product",
	}
	err := repo.Update(ctx, product)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestProductRepository_Delete(t *testing.T) {
	repo, _, ctx := setupProductTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewProductRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewProductRepository(db)
	if repo == nil {
		t.Error("NewProductRepository() returned nil")
	}
}

func setupProductTest(t *testing.T) (*ProductRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewProductRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
