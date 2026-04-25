package service

import (
	"context"

	"warehouse/internal/model"
)

type PermissionRepository interface {
	List(ctx context.Context, page, pageSize int) ([]model.Permission, int, error)
	GetByID(ctx context.Context, id int64) (*model.Permission, error)
}

type PermissionService struct {
	permRepo PermissionRepository
}

func NewPermissionService(permRepo PermissionRepository) *PermissionService {
	return &PermissionService{
		permRepo: permRepo,
	}
}

func (s *PermissionService) List(ctx context.Context, page, pageSize int) ([]model.Permission, int, error) {
	return s.permRepo.List(ctx, page, pageSize)
}

func (s *PermissionService) GetByID(ctx context.Context, id int64) (*model.Permission, error) {
	return s.permRepo.GetByID(ctx, id)
}
