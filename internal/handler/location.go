package handler

import (
	"context"
	"strconv"

	"warehouse/internal/model"
	"warehouse/internal/pkg/response"
	apperrors "warehouse/internal/pkg/errors"
	"warehouse/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateLocationRequest struct {
	WarehouseID int64  `json:"warehouse_id" binding:"required"`
	Zone        string `json:"zone" binding:"required"`
	Shelf       string `json:"shelf" binding:"required"`
	Level       string `json:"level" binding:"required"`
	Position    string `json:"position" binding:"required"`
	Status      *int   `json:"status"`
}

type UpdateLocationRequest struct {
	Zone     string `json:"zone"`
	Shelf    string `json:"shelf"`
	Level    string `json:"level"`
	Position string `json:"position"`
	Status   *int   `json:"status"`
}

type LocationListResponse struct {
	Locations []model.Location `json:"items"`
	Total     int              `json:"total"`
	Page      int              `json:"page"`
	Size      int              `json:"size"`
}

type LocationService interface {
	Create(ctx context.Context, input *service.CreateLocationInput) (*model.Location, error)
	GetByID(ctx context.Context, id int64) (*model.Location, error)
	List(ctx context.Context, page, pageSize int, warehouseID int64) (*service.ListLocationsResult, error)
	Update(ctx context.Context, id int64, input *service.UpdateLocationInput) (*model.Location, error)
	Delete(ctx context.Context, id int64) error
}

type LocationHandler struct {
	locationService LocationService
}

func NewLocationHandler(locationService LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

func (h *LocationHandler) Create(c *gin.Context) {
	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.CreateLocationInput{
		WarehouseID: req.WarehouseID,
		Zone:        req.Zone,
		Shelf:       req.Shelf,
		Level:       req.Level,
		Position:    req.Position,
	}

	if req.Status != nil {
		input.Status = *req.Status
	}

	location, err := h.locationService.Create(c.Request.Context(), input)
	if err != nil {
		handleLocationError(c, err)
		return
	}

	response.Success(c, location)
}

func (h *LocationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid location id")
		return
	}

	location, err := h.locationService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleLocationError(c, err)
		return
	}

	response.Success(c, location)
}

func (h *LocationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	warehouseID, _ := strconv.ParseInt(c.Query("warehouse_id"), 10, 64)

	result, err := h.locationService.List(c.Request.Context(), page, pageSize, warehouseID)
	if err != nil {
		handleLocationError(c, err)
		return
	}

	response.Success(c, LocationListResponse{
		Locations: result.Locations,
		Total:     result.Total,
		Page:      page,
		Size:      pageSize,
	})
}

func (h *LocationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid location id")
		return
	}

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid request body")
		return
	}

	input := &service.UpdateLocationInput{
		Zone:     req.Zone,
		Shelf:    req.Shelf,
		Level:    req.Level,
		Position: req.Position,
		Status:   req.Status,
	}

	location, err := h.locationService.Update(c.Request.Context(), id, input)
	if err != nil {
		handleLocationError(c, err)
		return
	}

	response.Success(c, location)
}

func (h *LocationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.CodeBadRequest, "invalid location id")
		return
	}

	err = h.locationService.Delete(c.Request.Context(), id)
	if err != nil {
		handleLocationError(c, err)
		return
	}

	response.Success(c, nil)
}

func handleLocationError(c *gin.Context, err error) {
	if appErr, ok := apperrors.GetAppError(err); ok {
		response.Error(c, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, apperrors.CodeInternalError, "internal server error")
}

func RegisterLocationRoutes(r *gin.RouterGroup, h *LocationHandler) {
	locations := r.Group("/locations")
	{
		locations.GET("", h.List)
		locations.POST("", h.Create)
		locations.GET("/:id", h.GetByID)
		locations.PUT("/:id", h.Update)
		locations.DELETE("/:id", h.Delete)
	}
}
