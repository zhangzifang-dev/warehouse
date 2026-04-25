package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type WarehouseRepository struct {
	db *bun.DB
}

func NewWarehouseRepository(db *bun.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) Create(ctx context.Context, warehouse *model.Warehouse) error {
	_, err := r.db.NewInsert().Model(warehouse).Exec(ctx)
	return err
}

func (r *WarehouseRepository) GetByID(ctx context.Context, id int64) (*model.Warehouse, error) {
	warehouse := new(model.Warehouse)
	err := r.db.NewSelect().
		Model(warehouse).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return warehouse, nil
}

func (r *WarehouseRepository) GetByCode(ctx context.Context, code string) (*model.Warehouse, error) {
	warehouse := new(model.Warehouse)
	err := r.db.NewSelect().
		Model(warehouse).
		Where("code = ?", code).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return warehouse, nil
}

func (r *WarehouseRepository) List(ctx context.Context, page, pageSize int) ([]model.Warehouse, int, error) {
	var warehouses []model.Warehouse
	total, err := r.db.NewSelect().
		Model(&warehouses).
		Where("deleted_at = ?", timeZero).
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return warehouses, total, nil
}

func (r *WarehouseRepository) Update(ctx context.Context, warehouse *model.Warehouse) error {
	_, err := r.db.NewUpdate().
		Model(warehouse).
		WherePK().
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *WarehouseRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Warehouse)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}


