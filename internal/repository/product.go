package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type ProductRepository struct {
	db *bun.DB
}

func NewProductRepository(db *bun.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	_, err := r.db.NewInsert().Model(product).Exec(ctx)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	product := new(model.Product)
	err := r.db.NewSelect().
		Model(product).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	product := new(model.Product)
	err := r.db.NewSelect().
		Model(product).
		Where("sku = ?", sku).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) List(ctx context.Context, page, pageSize int, categoryID int64, keyword string) ([]model.Product, int, error) {
	var products []model.Product
	query := r.db.NewSelect().
		Model(&products).
		Where("deleted_at = ?", timeZero)

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if keyword != "" {
		query = query.Where("name LIKE ? OR sku LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	total, err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	_, err := r.db.NewUpdate().
		Model(product).
		WherePK().
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Product)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}
