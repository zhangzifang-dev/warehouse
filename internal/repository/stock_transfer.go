package repository

import (
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type StockTransferRepository struct {
	db *bun.DB
}

func NewStockTransferRepository(db *bun.DB) *StockTransferRepository {
	return &StockTransferRepository{db: db}
}

func (r *StockTransferRepository) Create(ctx context.Context, transfer *model.StockTransfer) error {
	_, err := r.db.NewInsert().Model(transfer).Exec(ctx)
	return err
}

func (r *StockTransferRepository) GetByID(ctx context.Context, id int64) (*model.StockTransfer, error) {
	transfer := new(model.StockTransfer)
	err := r.db.NewSelect().
		Model(transfer).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r *StockTransferRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.StockTransfer, error) {
	transfer := new(model.StockTransfer)
	err := r.db.NewSelect().
		Model(transfer).
		Where("order_no = ?", orderNo).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r *StockTransferRepository) List(ctx context.Context, page, pageSize int, fromWarehouseID, toWarehouseID, status int) ([]model.StockTransfer, int, error) {
	var transfers []model.StockTransfer
	query := r.db.NewSelect().
		Model(&transfers).
		Where("deleted_at IS NULL")

	if fromWarehouseID > 0 {
		query = query.Where("source_warehouse_id = ?", fromWarehouseID)
	}
	if toWarehouseID > 0 {
		query = query.Where("target_warehouse_id = ?", toWarehouseID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	total, err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return transfers, total, nil
}

func (r *StockTransferRepository) ListWithFilter(ctx context.Context, filter *model.StockTransferQueryFilter) ([]model.StockTransfer, int, error) {
	var transfers []model.StockTransfer
	q := r.db.NewSelect().
		Model(&transfers).
		Where("deleted_at IS NULL")

	if filter.OrderNo != "" {
		q = q.Where("order_no LIKE ?", "%"+filter.OrderNo+"%")
	}

	if filter.SourceWarehouseID != nil {
		q = q.Where("source_warehouse_id = ?", *filter.SourceWarehouseID)
	}

	if filter.TargetWarehouseID != nil {
		q = q.Where("target_warehouse_id = ?", *filter.TargetWarehouseID)
	}

	if filter.CreatedAtStart != nil {
		q = q.Where("created_at >= ?", filter.CreatedAtStart)
	}

	if filter.CreatedAtEnd != nil {
		q = q.Where("created_at <= ?", filter.CreatedAtEnd)
	}

	total, err := q.
		Order("id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return transfers, total, nil
}

func (r *StockTransferRepository) Update(ctx context.Context, transfer *model.StockTransfer) error {
	_, err := r.db.NewUpdate().
		Model(transfer).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *StockTransferRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.StockTransfer)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
