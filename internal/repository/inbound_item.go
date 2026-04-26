package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type InboundItemRepository struct {
	db *bun.DB
}

func NewInboundItemRepository(db *bun.DB) *InboundItemRepository {
	return &InboundItemRepository{db: db}
}

func (r *InboundItemRepository) Create(ctx context.Context, item *model.InboundItem) error {
	_, err := r.db.NewInsert().Model(item).Exec(ctx)
	return err
}

func (r *InboundItemRepository) ListByOrderID(ctx context.Context, orderID int64) ([]model.InboundItem, error) {
	var items []model.InboundItem
	err := r.db.NewSelect().
		Model(&items).
		Where("order_id = ?", orderID).
		Where("deleted_at IS NULL").
		Order("id ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *InboundItemRepository) Update(ctx context.Context, item *model.InboundItem) error {
	_, err := r.db.NewUpdate().
		Model(item).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *InboundItemRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.InboundItem)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
