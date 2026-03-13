package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type TrackingService struct {
	tracking repository.TrackingRepository
}

func NewTrackingService(tracking repository.TrackingRepository) *TrackingService {
	return &TrackingService{tracking: tracking}
}

func (s *TrackingService) LogFood(ctx context.Context, userID int64, req domain.FoodLogRequest) (*domain.FoodLog, error) {
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	log := &domain.FoodLog{
		UserID:   userID,
		Date:     date,
		MealType: req.MealType,
		FoodName: req.FoodName,
		Barcode:  req.Barcode,
		Servings: req.Servings,
		Calories: req.Calories,
		ProteinG: req.ProteinG,
		CarbsG:   req.CarbsG,
		FatG:     req.FatG,
		FiberG:   req.FiberG,
	}

	if err := s.tracking.CreateFoodLog(ctx, log); err != nil {
		return nil, fmt.Errorf("logging food: %w", err)
	}
	return log, nil
}

func (s *TrackingService) GetFoodLogs(ctx context.Context, userID int64, date string) ([]domain.FoodLog, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	logs, err := s.tracking.GetFoodLogs(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("getting food logs: %w", err)
	}
	return logs, nil
}

func (s *TrackingService) DeleteFoodLog(ctx context.Context, id, userID int64) error {
	if err := s.tracking.DeleteFoodLog(ctx, id, userID); err != nil {
		return fmt.Errorf("deleting food log: %w", err)
	}
	return nil
}

func (s *TrackingService) LogWater(ctx context.Context, userID int64, req domain.WaterLogRequest) (*domain.WaterLog, error) {
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	log := &domain.WaterLog{
		UserID:   userID,
		Date:     date,
		AmountML: req.AmountML,
	}

	if err := s.tracking.CreateWaterLog(ctx, log); err != nil {
		return nil, fmt.Errorf("logging water: %w", err)
	}
	return log, nil
}

func (s *TrackingService) GetWaterLogs(ctx context.Context, userID int64, date string) ([]domain.WaterLog, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	logs, err := s.tracking.GetWaterLogs(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("getting water logs: %w", err)
	}
	return logs, nil
}

func (s *TrackingService) GetNutritionSummary(ctx context.Context, userID int64, date string) (*domain.NutritionSummary, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	summary, err := s.tracking.GetNutritionSummary(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("getting nutrition summary: %w", err)
	}
	return summary, nil
}
