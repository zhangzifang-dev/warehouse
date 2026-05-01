package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditLogService interface {
	GetByID(ctx context.Context, id int64) (*model.AuditLog, error)
	List(ctx context.Context, filter *service.AuditLogQueryFilter) (*service.AuditLogListResult, error)
}

type AuditLogHandler struct {
	auditLogService AuditLogService
}

func NewAuditLogHandler(auditLogService AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

type AuditLogListResponse struct {
	Items []model.AuditLog `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

func (h *AuditLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	filter := &service.AuditLogQueryFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if tableName := c.Query("table_name"); tableName != "" {
		filter.TableName = tableName
	}

	if recordIDStr := c.Query("record_id"); recordIDStr != "" {
		if recordID, err := strconv.ParseInt(recordIDStr, 10, 64); err == nil {
			filter.RecordID = &recordID
		}
	}

	if operatedByStr := c.Query("operated_by"); operatedByStr != "" {
		if operatedBy, err := strconv.ParseInt(operatedByStr, 10, 64); err == nil {
			filter.OperatedBy = &operatedBy
		}
	}

	if operatedByName := c.Query("operated_by_name"); operatedByName != "" {
		filter.OperatedByName = operatedByName
	}

	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = &startTime
		}
	}

	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = &endTime
		}
	}

	result, err := h.auditLogService.List(c.Request.Context(), filter)
	if err != nil {
		log.Printf("ERROR: failed to list audit logs: %v", err)
		response.Error(c, apperrors.CodeInternalError, "failed to list audit logs")
		return
	}

	response.Success(c, AuditLogListResponse{
		Items: result.Items,
		Total: result.Total,
		Page:  page,
		Size:  pageSize,
	})
}

func (h *AuditLogHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid audit log id")
		return
	}

	auditLog, err := h.auditLogService.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("ERROR: failed to get audit log: %v", err)
		response.Error(c, apperrors.CodeNotFound, "audit log not found")
		return
	}

	response.Success(c, auditLog)
}

func RegisterAuditLogRoutes(r *gin.RouterGroup, h *AuditLogHandler) {
	auditLogs := r.Group("/audit-logs")
	{
		auditLogs.GET("", h.List)
		auditLogs.GET("/:id", h.GetByID)
	}
}
