package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"warehouse/internal/model"
)

func TestCustomerRepository_Create(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	customer := &model.Customer{
		Name:   "Test Customer",
		Code:   "CUS001",
		Status: model.CustomerStatusActive,
	}

	err := repo.Create(ctx, customer)
	if err == nil {
		t.Error("Create() should return error with mock DB")
	}
}

func TestCustomerRepository_GetByID_Query(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	_, err := repo.GetByID(ctx, 1)
	if err == nil {
		t.Error("GetByID() should return error with mock DB")
	}
}

func TestCustomerRepository_GetByCode_Query(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	_, err := repo.GetByCode(ctx, "CUS001")
	if err == nil {
		t.Error("GetByCode() should return error with mock DB")
	}
}

func TestCustomerRepository_List_WithFilter(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	filter := &CustomerQueryFilter{
		Code:     "CUS",
		Name:     "Test",
		Phone:    "123",
		Status:   intPtrCustomer(1),
		Page:     1,
		PageSize: 10,
	}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestCustomerRepository_List_WithFilter_PartialFields(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	filter := &CustomerQueryFilter{
		Name:     "Test",
		Page:     1,
		PageSize: 10,
	}
	_, _, err := repo.List(ctx, filter)
	if err == nil {
		t.Error("List() should return error with mock DB")
	}
}

func TestCustomerRepository_Update_Query(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	customer := &model.Customer{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Updated Customer",
	}
	err := repo.Update(ctx, customer)
	if err == nil {
		t.Error("Update() should return error with mock DB")
	}
}

func TestCustomerRepository_Delete_Query(t *testing.T) {
	repo, _, ctx := setupCustomerTest(t)
	err := repo.Delete(ctx, 1)
	if err == nil {
		t.Error("Delete() should return error with mock DB")
	}
}

func TestNewCustomerRepository(t *testing.T) {
	sqldb, _ := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewCustomerRepository(db)
	if repo == nil {
		t.Error("NewCustomerRepository() returned nil")
	}
}

func setupCustomerTest(t *testing.T) (*CustomerRepository, *bun.DB, context.Context) {
	t.Helper()
	sqldb, err := sql.Open("mysql", "invalid:invalid@tcp(localhost:3306)/invalid")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	db := bun.NewDB(sqldb, mysqldialect.New())
	repo := NewCustomerRepository(db)
	ctx := context.Background()
	return repo, db, ctx
}

func intPtrCustomer(i int) *int {
	return &i
}
