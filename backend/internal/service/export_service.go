package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type ExportService struct {
	trackingRepo repository.TrackingRepository
	userRepo     repository.UserRepository
}

func NewExportService(trackingRepo repository.TrackingRepository, userRepo repository.UserRepository) *ExportService {
	return &ExportService{
		trackingRepo: trackingRepo,
		userRepo:     userRepo,
	}
}

func (s *ExportService) ExportFoodLogsCSV(ctx context.Context, userID int64, w io.Writer) error {
	logs, err := s.trackingRepo.GetFoodLogsByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("get food logs: %w", err)
	}

	cw := csv.NewWriter(w)
	defer cw.Flush()

	// Header
	if err := cw.Write([]string{"date", "meal_type", "food_name", "calories", "protein_g", "carbs_g", "fat_g", "fiber_g", "servings"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, l := range logs {
		if err := cw.Write([]string{
			l.Date,
			l.MealType,
			l.FoodName,
			strconv.Itoa(l.Calories),
			fmt.Sprintf("%.1f", l.ProteinG),
			fmt.Sprintf("%.1f", l.CarbsG),
			fmt.Sprintf("%.1f", l.FatG),
			fmt.Sprintf("%.1f", l.FiberG),
			fmt.Sprintf("%.1f", l.Servings),
		}); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	return nil
}

type FullExport struct {
	User     interface{} `json:"user"`
	FoodLogs interface{} `json:"food_logs"`
}

func (s *ExportService) ExportJSON(ctx context.Context, userID int64, w io.Writer) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	// Clear sensitive data
	user.PasswordHash = ""

	logs, err := s.trackingRepo.GetFoodLogsByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("get food logs: %w", err)
	}

	export := FullExport{
		User:     user,
		FoodLogs: logs,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(export)
}

func (s *ExportService) DeleteAccount(ctx context.Context, userID int64) error {
	return s.userRepo.Delete(ctx, userID)
}
