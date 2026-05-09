package service

import (
	"context"
	"log"
	"time"

	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/model"
)

type DashboardRepository interface {
	GetOverviewStats(ctx context.Context) (*model.OverviewStats, error)
	GetTrendData(ctx context.Context, params *model.DashboardQueryParams) ([]model.TrendData, error)
	GetTopProducts(ctx context.Context, params *model.DashboardQueryParams) ([]model.TopProduct, error)
	GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error)
	GetSupplierPerformance(ctx context.Context, params *model.DashboardQueryParams) ([]model.SupplierPerformance, error)
	GetPendingOrders(ctx context.Context) (*model.PendingOrders, error)
}

type DashboardService struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetOverview(ctx context.Context) (*model.OverviewStats, error) {
	stats, err := s.repo.GetOverviewStats(ctx)
	if err != nil {
		log.Printf("[ERROR] Service GetOverview: %v", err)
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get overview stats")
	}
	return stats, nil
}

func (s *DashboardService) GetTrendData(ctx context.Context, startDate, endDate time.Time) ([]model.TrendData, error) {
	if startDate.After(endDate) {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "start date must be before end date")
	}

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
	}

	trend, err := s.repo.GetTrendData(ctx, params)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get trend data")
	}

	return trend, nil
}

func (s *DashboardService) GetTopProducts(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.TopProduct, error) {
	if startDate.After(endDate) {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "start date must be before end date")
	}

	if limit <= 0 {
		limit = 10
	}

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	products, err := s.repo.GetTopProducts(ctx, params)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get top products")
	}

	return products, nil
}

func (s *DashboardService) GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error) {
	usage, err := s.repo.GetWarehouseUsage(ctx)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get warehouse usage")
	}

	return usage, nil
}

func (s *DashboardService) GetSupplierPerformance(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.SupplierPerformance, error) {
	if startDate.After(endDate) {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "start date must be before end date")
	}

	if limit <= 0 {
		limit = 10
	}

	params := &model.DashboardQueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	performance, err := s.repo.GetSupplierPerformance(ctx, params)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get supplier performance")
	}

	return performance, nil
}

func (s *DashboardService) GetPendingOrders(ctx context.Context) (*model.PendingOrders, error) {
	pending, err := s.repo.GetPendingOrders(ctx)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to get pending orders")
	}

	return pending, nil
}
