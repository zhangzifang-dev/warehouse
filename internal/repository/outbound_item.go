package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type OutboundItemRepository struct {
	db *bun.DB
}

func NewOutboundItemRepository(db *bun.DB) *OutboundItemRepository {
	return &OutboundItemRepository{db: db}
}

func (r *OutboundItemRepository) Create(ctx context.Context, item *model.OutboundItem) error {
	_, err := r.db.NewInsert().Model(item).Exec(ctx)
	return err
}

func (r *OutboundItemRepository) ListByOrderID(ctx context.Context, orderID int64) ([]model.OutboundItem, error) {
	var items []model.OutboundItem
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

func (r *OutboundItemRepository) Update(ctx context.Context, item *model.OutboundItem) error {
	_, err := r.db.NewUpdate().
		Model(item).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *OutboundItemRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.OutboundItem)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
