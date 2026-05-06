package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type OutboundOrderRepository struct {
	db *bun.DB
}

func NewOutboundOrderRepository(db *bun.DB) *OutboundOrderRepository {
	return &OutboundOrderRepository{db: db}
}

func (r *OutboundOrderRepository) Create(ctx context.Context, order *model.OutboundOrder) error {
	_, err := r.db.NewInsert().Model(order).Exec(ctx)
	return err
}

func (r *OutboundOrderRepository) GetByID(ctx context.Context, id int64) (*model.OutboundOrder, error) {
	order := new(model.OutboundOrder)
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

func (r *OutboundOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.OutboundOrder, error) {
	order := new(model.OutboundOrder)
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

func (r *OutboundOrderRepository) List(ctx context.Context, page, pageSize int, warehouseID, status int) ([]model.OutboundOrder, int, error) {
	var orders []model.OutboundOrder
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

func (r *OutboundOrderRepository) Update(ctx context.Context, order *model.OutboundOrder) error {
	_, err := r.db.NewUpdate().
		Model(order).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *OutboundOrderRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.OutboundOrder)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *OutboundOrderRepository) ListWithFilter(ctx context.Context, filter *model.OutboundOrderQueryFilter) ([]model.OutboundOrder, int, error) {
	var orders []model.OutboundOrder
	q := r.db.NewSelect().
		Model(&orders).
		Where("deleted_at IS NULL")

	if filter.OrderNo != "" {
		q = q.Where("order_no LIKE ?", "%"+filter.OrderNo+"%")
	}

	if filter.CustomerID != nil {
		q = q.Where("customer_id = ?", *filter.CustomerID)
	}

	if filter.WarehouseID != nil {
		q = q.Where("warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.QuantityMin != nil {
		q = q.Where("total_quantity >= ?", *filter.QuantityMin)
	}

	if filter.QuantityMax != nil {
		q = q.Where("total_quantity <= ?", *filter.QuantityMax)
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
	return orders, total, nil
}
