package repository

import (
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]model.User, int, error) {
	var users []model.User
	total, err := r.db.NewSelect().
		Model(&users).
		Where("deleted_at IS NULL").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.User)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *UserRepository) GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.NewSelect().
		Model(&roles).
		Join("JOIN user_roles ur ON ur.role_id = role.id").
		Where("ur.user_id = ?", userID).
		Where("ur.deleted_at IS NULL").
		Where("role.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) GetUserPermissions(ctx context.Context, userID int64) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.NewSelect().
		Model(&permissions).
		Join("JOIN role_permissions rp ON rp.permission_id = permission.id").
		Join("JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ?", userID).
		Where("ur.deleted_at IS NULL").
		Where("rp.deleted_at IS NULL").
		Where("permission.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *UserRepository) AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	userRoles := make([]model.UserRole, len(roleIDs))
	for i, roleID := range roleIDs {
		userRoles[i] = model.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
	}

	_, err := r.db.NewInsert().Model(&userRoles).Exec(ctx)
	return err
}

func (r *UserRepository) CreateRole(ctx context.Context, role *model.Role) error {
	_, err := r.db.NewInsert().Model(role).Exec(ctx)
	return err
}

func (r *UserRepository) CreatePermission(ctx context.Context, perm *model.Permission) error {
	_, err := r.db.NewInsert().Model(perm).Exec(ctx)
	return err
}

func (r *UserRepository) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
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
