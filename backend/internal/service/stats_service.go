package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type StatsService struct {
	stats repository.StatsRepository
}

func NewStatsService(stats repository.StatsRepository) *StatsService {
	return &StatsService{stats: stats}
}

func (s *StatsService) GetWeeklyStats(ctx context.Context, userID int64) (*domain.WeeklyStats, error) {
	now := time.Now().UTC()
	// Go back to Monday of this week
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)

	startDate := monday.Format("2006-01-02")
	endDate := sunday.Format("2006-01-02")

	days, err := s.stats.GetWeeklyStats(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("getting weekly stats: %w", err)
	}

	result := &domain.WeeklyStats{Days: days}
	if len(days) > 0 {
		totalCals := 0
		totalWater := 0
		for _, d := range days {
			totalCals += d.Calories
			totalWater += d.WaterML
		}
		result.AvgCals = totalCals / len(days)
		result.AvgWater = totalWater / len(days)
	}

	return result, nil
}

func (s *StatsService) GetMonthlyStats(ctx context.Context, userID int64) (*domain.MonthlyStats, error) {
	now := time.Now().UTC()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	days, err := s.stats.GetMonthlyStats(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("getting monthly stats: %w", err)
	}

	result := &domain.MonthlyStats{Days: days}
	if len(days) > 0 {
		totalCals := 0
		totalWater := 0
		for _, d := range days {
			totalCals += d.Calories
			totalWater += d.WaterML
		}
		result.AvgCals = totalCals / len(days)
		result.AvgWater = totalWater / len(days)
	}

	return result, nil
}

func (s *StatsService) GetStreaks(ctx context.Context, userID int64) (*domain.StreakInfo, error) {
	streaks, err := s.stats.GetStreaks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting streaks: %w", err)
	}
	return streaks, nil
}
