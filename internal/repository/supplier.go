package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type SupplierRepository struct {
	db *bun.DB
}

func NewSupplierRepository(db *bun.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
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
		Where("deleted_at = ?", timeZero).
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
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (r *SupplierRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]model.Supplier, int, error) {
	var suppliers []model.Supplier
	query := r.db.NewSelect().
		Model(&suppliers).
		Where("deleted_at = ?", timeZero)

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	total, err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
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
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *SupplierRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Supplier)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}
