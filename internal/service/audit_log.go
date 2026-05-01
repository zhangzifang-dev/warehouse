package service

import (
	"context"
	"time"

	"warehouse/internal/model"
	"warehouse/internal/repository"
)

type AuditLogFullRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	GetByID(ctx context.Context, id int64) (*model.AuditLog, error)
	List(ctx context.Context, filter *repository.AuditLogFilter) ([]model.AuditLog, int, error)
}

type AuditLogService struct {
	repo AuditLogFullRepository
}

func NewAuditLogService(repo AuditLogFullRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

type CreateAuditLogInput struct {
	TableName  string
	RecordID   int64
	Action     string
	OldValue   map[string]any
	NewValue   map[string]any
	OperatedBy int64
	IPAddress  string
}

func (s *AuditLogService) Log(ctx context.Context, input *CreateAuditLogInput) error {
	log := &model.AuditLog{
		TableName:  input.TableName,
		RecordID:   input.RecordID,
		Action:     input.Action,
		OldValue:   input.OldValue,
		NewValue:   input.NewValue,
		OperatedBy: input.OperatedBy,
		OperatedAt: time.Now(),
		IPAddress:  input.IPAddress,
	}

	err := s.repo.Create(ctx, log)
	if err != nil {
	}
	return err
}

type AuditLogQueryFilter struct {
	TableName       string
	RecordID        *int64
	OperatedBy      *int64
	OperatedByName  string
	Action          string
	StartTime       *time.Time
	EndTime         *time.Time
	Page            int
	PageSize        int
}

type AuditLogListResult struct {
	Items []model.AuditLog
	Total int
}

func (s *AuditLogService) GetByID(ctx context.Context, id int64) (*model.AuditLog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AuditLogService) List(ctx context.Context, filter *AuditLogQueryFilter) (*AuditLogListResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	repoFilter := &repository.AuditLogFilter{
		TableName:  filter.TableName,
		RecordID:   filter.RecordID,
		OperatedBy: filter.OperatedBy,
		OperatedByName:  filter.OperatedByName,
		Action:         filter.Action,
		StartTime:  filter.StartTime,
		EndTime:    filter.EndTime,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}

	items, total, err := s.repo.List(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	return &AuditLogListResult{
		Items: items,
		Total: total,
	}, nil
}
