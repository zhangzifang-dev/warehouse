package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/repository"
)

type CategoryQueryFilter struct {
	Name     string
	Page     int
	PageSize int
}

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	List(ctx context.Context, filter *repository.CategoryQueryFilter) ([]model.Category, int, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	HasChildren(ctx context.Context, id int64) (bool, error)
}

type CategoryService struct {
	categoryRepo CategoryRepository
	auditLogger  AuditLogger
}

func NewCategoryService(categoryRepo CategoryRepository, auditLogger AuditLogger) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		auditLogger:  auditLogger,
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

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(category)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "categories",
			RecordID:   category.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: category.CreatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
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

func (s *CategoryService) List(ctx context.Context, filter *CategoryQueryFilter) (*ListCategoriesResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	repoFilter := &repository.CategoryQueryFilter{
		Name:     filter.Name,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}

	categories, total, err := s.categoryRepo.List(ctx, repoFilter)
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

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(category)
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

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(category)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "categories",
			RecordID:   category.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: category.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
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

	if s.auditLogger != nil {
		oldValue, _ := json.Marshal(category)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "categories",
			RecordID:   category.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: category.UpdatedBy,
			IPAddress:  GetClientIPFromContext(ctx),
		})
	}

	return nil
}
