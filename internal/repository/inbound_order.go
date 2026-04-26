package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type InboundOrderRepository struct {
	db *bun.DB
}

func NewInboundOrderRepository(db *bun.DB) *InboundOrderRepository {
	return &InboundOrderRepository{db: db}
}

func (r *InboundOrderRepository) Create(ctx context.Context, order *model.InboundOrder) error {
	_, err := r.db.NewInsert().Model(order).Exec(ctx)
	return err
}

func (r *InboundOrderRepository) GetByID(ctx context.Context, id int64) (*model.InboundOrder, error) {
	order := new(model.InboundOrder)
	err := r.db.NewSelect().
		Model(order).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *InboundOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.InboundOrder, error) {
	order := new(model.InboundOrder)
	err := r.db.NewSelect().
		Model(order).
		Where("order_no = ?", orderNo).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *InboundOrderRepository) List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.InboundOrder, int, error) {
	var orders []model.InboundOrder
	query := r.db.NewSelect().
		Model(&orders).
		Where("deleted_at IS NULL")

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if status > 0 {
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
	return orders, total, nil
}

func (r *InboundOrderRepository) Update(ctx context.Context, order *model.InboundOrder) error {
	_, err := r.db.NewUpdate().
		Model(order).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *InboundOrderRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.InboundOrder)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
