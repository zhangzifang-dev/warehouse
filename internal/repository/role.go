package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type RoleRepository struct {
	db *bun.DB
}

func NewRoleRepository(db *bun.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *model.Role) error {
	_, err := r.db.NewInsert().Model(role).Exec(ctx)
	return err
}

func (r *RoleRepository) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	role := new(model.Role)
	err := r.db.NewSelect().
		Model(role).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	role := new(model.Role)
	err := r.db.NewSelect().
		Model(role).
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) List(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
	var roles []model.Role
	total, err := r.db.NewSelect().
		Model(&roles).
		Where("deleted_at IS NULL").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *model.Role) error {
	_, err := r.db.NewUpdate().
		Model(role).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Role)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *RoleRepository) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	rolePerms := make([]model.RolePermission, len(permissionIDs))
	for i, permID := range permissionIDs {
		rolePerms[i] = model.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
	}

	_, err := r.db.NewInsert().Model(&rolePerms).Exec(ctx)
	return err
}

func (r *RoleRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.NewSelect().
		Model(&permissions).
		Join("JOIN role_permissions rp ON rp.permission_id = permission.id").
		Where("rp.role_id = ?", roleID).
		Where("rp.deleted_at IS NULL").
		Where("permission.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
