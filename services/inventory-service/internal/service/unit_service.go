package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/erp-system/services/inventory-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitService struct {
	unitRepo      *repository.UnitRepository
	unitChartRepo *repository.UnitChartRepository
}

func NewUnitService(
	unitRepo *repository.UnitRepository,
	unitChartRepo *repository.UnitChartRepository,
) *UnitService {
	return &UnitService{
		unitRepo:      unitRepo,
		unitChartRepo: unitChartRepo,
	}
}

type CreateUnitRequest struct {
	Code        string                 `json:"code" binding:"required"`
	Name        string                 `json:"name" binding:"required"`
	Symbol      string                 `json:"symbol"`
	Description string                 `json:"description"`
	UnitType    string                 `json:"unit_type" binding:"required"` // quantity, weight, volume, length
	IsBaseUnit  bool                   `json:"is_base_unit"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type UpdateUnitRequest struct {
	Name        *string                `json:"name"`
	Symbol      *string                `json:"symbol"`
	Description *string                `json:"description"`
	IsActive    *bool                  `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func (s *UnitService) CreateUnit(
	ctx context.Context,
	req CreateUnitRequest,
	createdBy primitive.ObjectID,
) (*models.Unit, error) {

	// Enforce single base unit per unit_type
	if req.IsBaseUnit {
		exists, err := s.unitRepo.ExistsBaseUnit(ctx, req.UnitType)
		if err != nil {
			return nil, fmt.Errorf("failed to check base unit: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("base unit already exists for unit type %s", req.UnitType)
		}
	}

	unit := &models.Unit{
		Code:        req.Code,
		Name:        req.Name,
		Symbol:      req.Symbol,
		Description: req.Description,
		UnitType:    req.UnitType,
		IsBaseUnit:  req.IsBaseUnit,
		IsActive:    true,
		Metadata:    req.Metadata,
	}

	unit.BaseModel.ID = primitive.NewObjectID()
	unit.BaseModel.CreatedAt = time.Now()
	unit.BaseModel.UpdatedAt = time.Now()
	unit.BaseModel.Version = 1
	unit.BaseModel.CreatedBy = createdBy

	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}

	return unit, nil
}

func (s *UnitService) DeleteUnit(
	ctx context.Context,
	unitID primitive.ObjectID,
) error {

	unit, err := s.unitRepo.FindByID(ctx, unitID)
	if err != nil {
		return fmt.Errorf("unit not found: %w", err)
	}

	if unit.IsBaseUnit {
		return fmt.Errorf("base unit cannot be deleted")
	}

	hasConversions, err := s.unitChartRepo.HasConversions(ctx, unitID)
	if err != nil {
		return fmt.Errorf("failed to check unit conversions: %w", err)
	}

	if hasConversions {
		return fmt.Errorf("unit has active conversions and cannot be deleted")
	}

	if err := s.unitRepo.SoftDelete(ctx, unitID); err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}

	return nil
}

func (s *UnitService) GetUnits(
	ctx context.Context,
	unitType *string,
	activeOnly bool,
) ([]*models.Unit, error) {

	units, err := s.unitRepo.Find(ctx, unitType, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get units: %w", err)
	}

	return units, nil
}

func (s *UnitService) GetUnit(
	ctx context.Context,
	unitID primitive.ObjectID,
) (*models.Unit, error) {

	unit, err := s.unitRepo.FindByID(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("unit not found: %w", err)
	}

	return unit, nil
}

func (s *UnitService) UpdateUnit(
	ctx context.Context,
	unitID primitive.ObjectID,
	req UpdateUnitRequest,
	updatedBy primitive.ObjectID,
) (*models.Unit, error) {

	unit, err := s.unitRepo.FindByID(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("unit not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		unit.Name = *req.Name
	}

	if req.Symbol != nil {
		unit.Symbol = *req.Symbol
	}

	if req.Description != nil {
		unit.Description = *req.Description
	}

	if req.IsActive != nil {
		unit.IsActive = *req.IsActive
	}

	if req.Metadata != nil {
		unit.Metadata = req.Metadata
	}

	unit.UpdatedBy = updatedBy
	unit.UpdatedAt = time.Now()
	unit.Version++

	if err := s.unitRepo.Update(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to update unit: %w", err)
	}

	return unit, nil
}
