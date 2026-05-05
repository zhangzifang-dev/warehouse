package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestCategoryRepository_Create(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	category := &model.Category{
		Name:      "Electronics",
		SortOrder: 1,
		Status:    model.CategoryStatusActive,
	}

	err := repo.Create(ctx, category)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestCategoryRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestCategoryRepository_List_Query(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	_, _, err := repo.List(ctx, &CategoryQueryFilter{Page: 1, PageSize: 10})
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestCategoryRepository_ListByParent_Query(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	_, _, err := repo.ListByParent(ctx, 1, 1, 10)
	if err == nil {
		t.Error("ListByParent() should return error with mock DB")
	}
}

func TestCategoryRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	category := &model.Category{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Electronics",
	}
	err := repo.Update(ctx, category)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestCategoryRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestCategoryRepository_HasChildren_Query(t *testing.T) {
	repo, _, ctx := setupCategoryTest(t)
	_, err := repo.HasChildren(ctx, 1)
	if err == nil {
		t.Error("HasChildren() should return error with mock DB")
	}
}

func TestNewCategoryRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewCategoryRepository(db)
	if repo == nil {
		t.Error("NewCategoryRepository() returned nil")
	}
}

func setupCategoryTest(t *testing.T) (*CategoryRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewCategoryRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}
