package service

import (
	"context"

	"warehouse/internal/model"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id int64) (*model.Role, error)
	GetByCode(ctx context.Context, code string) (*model.Role, error)
	List(ctx context.Context, page, pageSize int) ([]model.Role, int, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id int64) error
	AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error)
}

type RoleService struct {
	roleRepo RoleRepository
}

func NewRoleService(roleRepo RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

func (s *RoleService) Create(ctx context.Context, role *model.Role) error {
	return s.roleRepo.Create(ctx, role)
}

func (s *RoleService) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

func (s *RoleService) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	return s.roleRepo.GetByCode(ctx, code)
}

func (s *RoleService) List(ctx context.Context, page, pageSize int) ([]model.Role, int, error) {
	return s.roleRepo.List(ctx, page, pageSize)
}

func (s *RoleService) Update(ctx context.Context, role *model.Role) error {
	return s.roleRepo.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id int64) error {
	return s.roleRepo.Delete(ctx, id)
}

func (s *RoleService) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}

func (s *RoleService) GetRolePermissions(ctx context.Context, roleID int64) ([]model.Permission, error) {
	return s.roleRepo.GetRolePermissions(ctx, roleID)
}
