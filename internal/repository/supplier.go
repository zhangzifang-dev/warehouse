package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type SupplierRepository struct {
	db *bun.DB
}

func NewSupplierRepository(db *bun.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

type SupplierQueryFilter struct {
	Code     string
	Name     string
	Contact  string
	Phone    string
	Status   *int
	Page     int
	PageSize int
}

func (r *SupplierRepository) Create(ctx context.Context, supplier *model.Supplier) error {
	_, err := r.db.NewInsert().Model(supplier).Exec(ctx)
	return err
}

func (r *SupplierRepository) GetByID(ctx context.Context, id int64) (*model.Supplier, error) {
	supplier := new(model.Supplier)
	err := r.db.NewSelect().
		Model(supplier).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (r *SupplierRepository) GetByCode(ctx context.Context, code string) (*model.Supplier, error) {
	supplier := new(model.Supplier)
	err := r.db.NewSelect().
		Model(supplier).
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (r *SupplierRepository) List(ctx context.Context, filter *SupplierQueryFilter) ([]model.Supplier, int, error) {
	var suppliers []model.Supplier

	q := r.db.NewSelect().
		Model(&suppliers).
		Where("deleted_at IS NULL")

	if filter.Code != "" {
		q = q.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Contact != "" {
		q = q.Where("contact LIKE ?", "%"+filter.Contact+"%")
	}
	if filter.Phone != "" {
		q = q.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}

	total, err := q.
		Order("id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return suppliers, total, nil
}

func (r *SupplierRepository) Update(ctx context.Context, supplier *model.Supplier) error {
	_, err := r.db.NewUpdate().
		Model(supplier).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *SupplierRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Supplier)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
