package service

import (
	"context"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type LocationRepository interface {
	Create(ctx context.Context, location *model.Location) error
	GetByID(ctx context.Context, id int64) (*model.Location, error)
	GetByWarehouseAndCode(ctx context.Context, warehouseID int64, code string) (*model.Location, error)
	List(ctx context.Context, page, pageSize int, warehouseID int64) ([]model.Location, int, error)
	Update(ctx context.Context, location *model.Location) error
	Delete(ctx context.Context, id int64) error
}

type LocationService struct {
	locationRepo  LocationRepository
	warehouseRepo WarehouseRepository
}

func NewLocationService(locationRepo LocationRepository, warehouseRepo WarehouseRepository) *LocationService {
	return &LocationService{
		locationRepo:  locationRepo,
		warehouseRepo: warehouseRepo,
	}
}

type CreateLocationInput struct {
	WarehouseID int64
	Zone        string
	Shelf       string
	Level       string
	Position    string
	Status      int
}

type UpdateLocationInput struct {
	Zone     string
	Shelf    string
	Level    string
	Position string
	Status   *int
}

type ListLocationsResult struct {
	Locations []model.Location
	Total     int
}

func (s *LocationService) Create(ctx context.Context, input *CreateLocationInput) (*model.Location, error) {
	_, err := s.warehouseRepo.GetByID(ctx, input.WarehouseID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "warehouse not found")
	}

	code := input.Zone + "-" + input.Shelf + "-" + input.Level + "-" + input.Position

	existing, err := s.locationRepo.GetByWarehouseAndCode(ctx, input.WarehouseID, code)
	if err == nil && existing != nil {
		return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "location code already exists in this warehouse")
	}

	location := &model.Location{
		WarehouseID: input.WarehouseID,
		Zone:        input.Zone,
		Shelf:       input.Shelf,
		Level:       input.Level,
		Position:    input.Position,
		Code:        code,
		Status:      input.Status,
	}

	if location.Status == 0 {
		location.Status = model.LocationStatusActive
	}

	err = s.locationRepo.Create(ctx, location)
	if err != nil {
		if isDuplicateEntry(err) {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "location code already exists in this warehouse")
		}
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create location")
	}

	return location, nil
}

func (s *LocationService) GetByID(ctx context.Context, id int64) (*model.Location, error) {
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "location not found")
	}
	return location, nil
}

func (s *LocationService) List(ctx context.Context, page, pageSize int, warehouseID int64) (*ListLocationsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	locations, total, err := s.locationRepo.List(ctx, page, pageSize, warehouseID)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list locations")
	}

	return &ListLocationsResult{
		Locations: locations,
		Total:     total,
	}, nil
}

func (s *LocationService) Update(ctx context.Context, id int64, input *UpdateLocationInput) (*model.Location, error) {
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "location not found")
	}

	needCodeUpdate := false
	if input.Zone != "" {
		location.Zone = input.Zone
		needCodeUpdate = true
	}
	if input.Shelf != "" {
		location.Shelf = input.Shelf
		needCodeUpdate = true
	}
	if input.Level != "" {
		location.Level = input.Level
		needCodeUpdate = true
	}
	if input.Position != "" {
		location.Position = input.Position
		needCodeUpdate = true
	}
	if input.Status != nil {
		location.Status = *input.Status
	}

	if needCodeUpdate {
		newCode := location.Zone + "-" + location.Shelf + "-" + location.Level + "-" + location.Position
		existing, err := s.locationRepo.GetByWarehouseAndCode(ctx, location.WarehouseID, newCode)
		if err == nil && existing != nil && existing.ID != location.ID {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "location code already exists in this warehouse")
		}
		location.Code = newCode
	}

	err = s.locationRepo.Update(ctx, location)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update location")
	}

	return location, nil
}

func (s *LocationService) Delete(ctx context.Context, id int64) error {
	_, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "location not found")
	}

	err = s.locationRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete location")
	}

	return nil
}
