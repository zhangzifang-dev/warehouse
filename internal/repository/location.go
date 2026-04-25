package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type LocationRepository struct {
	db *bun.DB
}

func NewLocationRepository(db *bun.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(ctx context.Context, location *model.Location) error {
	_, err := r.db.NewInsert().Model(location).Exec(ctx)
	return err
}

func (r *LocationRepository) GetByID(ctx context.Context, id int64) (*model.Location, error) {
	location := new(model.Location)
	err := r.db.NewSelect().
		Model(location).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return location, nil
}

func (r *LocationRepository) GetByWarehouseAndCode(ctx context.Context, warehouseID int64, code string) (*model.Location, error) {
	location := new(model.Location)
	err := r.db.NewSelect().
		Model(location).
		Where("warehouse_id = ?", warehouseID).
		Where("code = ?", code).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return location, nil
}

func (r *LocationRepository) List(ctx context.Context, page, pageSize int, warehouseID int64) ([]model.Location, int, error) {
	var locations []model.Location
	query := r.db.NewSelect().
		Model(&locations).
		Where("deleted_at = ?", timeZero)

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	total, err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return locations, total, nil
}

func (r *LocationRepository) ListByWarehouse(ctx context.Context, warehouseID int64, page, pageSize int) ([]model.Location, int, error) {
	return r.List(ctx, page, pageSize, warehouseID)
}

func (r *LocationRepository) Update(ctx context.Context, location *model.Location) error {
	_, err := r.db.NewUpdate().
		Model(location).
		WherePK().
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *LocationRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Location)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}
