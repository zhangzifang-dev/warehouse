package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type CategoryRepository struct {
	db *bun.DB
}

func NewCategoryRepository(db *bun.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	_, err := r.db.NewInsert().Model(category).Exec(ctx)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	category := new(model.Category)
	err := r.db.NewSelect().
		Model(category).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) List(ctx context.Context, page, pageSize int, parentID int64) ([]model.Category, int, error) {
	var categories []model.Category
	query := r.db.NewSelect().
		Model(&categories).
		Where("deleted_at IS NULL")

	if parentID > 0 {
		query = query.Where("parent_id = ?", parentID)
	}

	total, err := query.
		Order("sort_order ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *CategoryRepository) ListByParent(ctx context.Context, parentID int64, page, pageSize int) ([]model.Category, int, error) {
	return r.List(ctx, page, pageSize, parentID)
}

func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	_, err := r.db.NewUpdate().
		Model(category).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Category)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *CategoryRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	count, err := r.db.NewSelect().
		Model((*model.Category)(nil)).
		Where("parent_id = ?", id).
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
