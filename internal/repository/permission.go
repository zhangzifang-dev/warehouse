package repository

import (
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type PermissionRepository struct {
	db *bun.DB
}

func NewPermissionRepository(db *bun.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, perm *model.Permission) error {
	_, err := r.db.NewInsert().Model(perm).Exec(ctx)
	return err
}

func (r *PermissionRepository) GetByID(ctx context.Context, id int64) (*model.Permission, error) {
	perm := new(model.Permission)
	err := r.db.NewSelect().
		Model(perm).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

func (r *PermissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	perm := new(model.Permission)
	err := r.db.NewSelect().
		Model(perm).
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

func (r *PermissionRepository) List(ctx context.Context, page, pageSize int) ([]model.Permission, int, error) {
	var perms []model.Permission
	total, err := r.db.NewSelect().
		Model(&perms).
		Where("deleted_at IS NULL").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return perms, total, nil
}

func (r *PermissionRepository) Update(ctx context.Context, perm *model.Permission) error {
	_, err := r.db.NewUpdate().
		Model(perm).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *PermissionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Permission)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *PermissionRepository) BatchCreate(ctx context.Context, perms []*model.Permission) error {
	if len(perms) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&perms).Exec(ctx)
	return err
}
