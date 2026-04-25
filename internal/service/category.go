package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	List(ctx context.Context, page, pageSize int, parentID int64) ([]model.Category, int, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	HasChildren(ctx context.Context, id int64) (bool, error)
}

type CategoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(categoryRepo CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

type CreateCategoryInput struct {
	Name      string
	ParentID  int64
	SortOrder int
	Status    int
}

type UpdateCategoryInput struct {
	Name      *string
	ParentID  *int64
	SortOrder *int
	Status    *int
}

type ListCategoriesResult struct {
	Categories []model.Category
	Total      int
}

func (s *CategoryService) Create(ctx context.Context, input *CreateCategoryInput) (*model.Category, error) {
	if input.Name == "" {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "category name is required")
	}

	category := &model.Category{
		Name:      input.Name,
		ParentID:  input.ParentID,
		SortOrder: input.SortOrder,
		Status:    input.Status,
	}

	if category.Status == 0 {
		category.Status = model.CategoryStatusActive
	}

	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create category")
	}

	return category, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "category not found")
	}
	return category, nil
}

func (s *CategoryService) List(ctx context.Context, page, pageSize int, parentID int64) (*ListCategoriesResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	categories, total, err := s.categoryRepo.List(ctx, page, pageSize, parentID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list categories")
	}

	return &ListCategoriesResult{
		Categories: categories,
		Total:      total,
	}, nil
}

func (s *CategoryService) Update(ctx context.Context, id int64, input *UpdateCategoryInput) (*model.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "category not found")
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.ParentID != nil {
		category.ParentID = *input.ParentID
	}
	if input.SortOrder != nil {
		category.SortOrder = *input.SortOrder
	}
	if input.Status != nil {
		category.Status = *input.Status
	}

	err = s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update category")
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "category not found")
	}

	hasChildren, err := s.categoryRepo.HasChildren(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to check children")
	}
	if hasChildren {
		return apperrors.NewAppError(apperrors.CodeBadRequest, "cannot delete category with children")
	}

	err = s.categoryRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete category")
	}

	return nil
}
