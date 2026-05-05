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

type CategoryQueryFilter struct {
	Name     string
	ParentID int64
	Page     int
	PageSize int
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

func (r *CategoryRepository) List(ctx context.Context, filter *CategoryQueryFilter) ([]model.Category, int, error) {
	var categories []model.Category

	q := r.db.NewSelect().
		Model(&categories).
		Where("deleted_at IS NULL")

	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.ParentID > 0 {
		q = q.Where("parent_id = ?", filter.ParentID)
	}

	total, err := q.
		Order("sort_order ASC, id ASC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *CategoryRepository) ListByParent(ctx context.Context, parentID int64, page, pageSize int) ([]model.Category, int, error) {
	filter := &CategoryQueryFilter{
		ParentID: parentID,
		Page:     page,
		PageSize: pageSize,
	}
	return r.List(ctx, filter)
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
