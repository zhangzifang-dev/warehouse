package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	GetOverview(ctx context.Context) (*model.OverviewStats, error)
	GetTrendData(ctx context.Context, startDate, endDate time.Time) ([]model.TrendData, error)
	GetTopProducts(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.TopProduct, error)
	GetWarehouseUsage(ctx context.Context) ([]model.WarehouseUsage, error)
	GetSupplierPerformance(ctx context.Context, startDate, endDate time.Time, limit int) ([]model.SupplierPerformance, error)
	GetPendingOrders(ctx context.Context) (*model.PendingOrders, error)
}

type DashboardHandler struct {
	dashboardService DashboardService
}

func NewDashboardHandler(dashboardService DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetOverview(c *gin.Context) {
	stats, err := h.dashboardService.GetOverview(c.Request.Context())
	if err != nil {
		log.Printf("[ERROR] GetOverview failed: %v", err)
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get overview stats")
		return
	}

	response.Success(c, stats)
}

func (h *DashboardHandler) GetTrendData(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid end_date format")
		return
	}

	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	trend, err := h.dashboardService.GetTrendData(c.Request.Context(), startDate, endDate)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get trend data")
		return
	}

	response.Success(c, trend)
}

func (h *DashboardHandler) GetTopProducts(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	limitStr := c.DefaultQuery("limit", "10")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid end_date format")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	products, err := h.dashboardService.GetTopProducts(c.Request.Context(), startDate, endDate, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get top products")
		return
	}

	response.Success(c, products)
}

func (h *DashboardHandler) GetWarehouseUsage(c *gin.Context) {
	usage, err := h.dashboardService.GetWarehouseUsage(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get warehouse usage")
		return
	}

	response.Success(c, usage)
}

func (h *DashboardHandler) GetSupplierPerformance(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	limitStr := c.DefaultQuery("limit", "10")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid end_date format")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	performance, err := h.dashboardService.GetSupplierPerformance(c.Request.Context(), startDate, endDate, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get supplier performance")
		return
	}

	response.Success(c, performance)
}

func (h *DashboardHandler) GetPendingOrders(c *gin.Context) {
	pending, err := h.dashboardService.GetPendingOrders(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, apperrors.CodeInternalError, "failed to get pending orders")
		return
	}

	response.Success(c, pending)
}
