package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type InventoryRepository struct {
	db *bun.DB
}

func NewInventoryRepository(db *bun.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(ctx context.Context, inventory *model.Inventory) error {
	_, err := r.db.NewInsert().Model(inventory).Exec(ctx)
	return err
}

func (r *InventoryRepository) GetByID(ctx context.Context, id int64) (*model.Inventory, error) {
	inventory := new(model.Inventory)
	err := r.db.NewSelect().
		Model(inventory).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *InventoryRepository) List(ctx context.Context, page, pageSize int, warehouseID, productID int64) ([]model.Inventory, int, error) {
	var inventories []model.Inventory
	query := r.db.NewSelect().
		Model(&inventories).
		Where("deleted_at = ?", timeZero)

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}

	total, err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return inventories, total, nil
}

func (r *InventoryRepository) Update(ctx context.Context, inventory *model.Inventory) error {
	_, err := r.db.NewUpdate().
		Model(inventory).
		WherePK().
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *InventoryRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Inventory)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *InventoryRepository) GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
	inventory := new(model.Inventory)
	query := r.db.NewSelect().
		Model(inventory).
		Where("warehouse_id = ?", warehouseID).
		Where("product_id = ?", productID).
		Where("deleted_at = ?", timeZero)

	if batchNo != "" {
		query = query.Where("batch_no = ?", batchNo)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *InventoryRepository) UpdateQuantity(ctx context.Context, id int64, quantity float64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Inventory)(nil)).
		Set("quantity = ?", quantity).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}
