package repository

import (
	"time"
	"context"

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
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *InventoryRepository) List(ctx context.Context, filter *model.InventoryQueryFilter) ([]model.Inventory, int, error) {
	var inventories []model.Inventory
	q := r.db.NewSelect().
		Model(&inventories).
		Relation("Product").
		Relation("Warehouse").
		Where("inventory.deleted_at IS NULL")

	if filter.ProductID > 0 {
		q = q.Where("inventory.product_id = ?", filter.ProductID)
	}

	if filter.WarehouseID > 0 {
		q = q.Where("inventory.warehouse_id = ?", filter.WarehouseID)
	}

	if filter.ProductName != "" {
		q = q.Where("product.name LIKE ?", "%"+filter.ProductName+"%")
	}

	if filter.QuantityMin != nil {
		q = q.Where("inventory.quantity >= ?", *filter.QuantityMin)
	}

	if filter.QuantityMax != nil {
		q = q.Where("inventory.quantity <= ?", *filter.QuantityMax)
	}

	if filter.BatchNo != "" {
		q = q.Where("inventory.batch_no LIKE ?", "%"+filter.BatchNo+"%")
	}

	total, err := q.
		Order("inventory.id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
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
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *InventoryRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Inventory)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *InventoryRepository) GetByWarehouseAndProduct(ctx context.Context, warehouseID, productID int64, batchNo string) (*model.Inventory, error) {
	inventory := new(model.Inventory)
	query := r.db.NewSelect().
		Model(inventory).
		Where("warehouse_id = ?", warehouseID).
		Where("product_id = ?", productID).
		Where("deleted_at IS NULL")

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
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
