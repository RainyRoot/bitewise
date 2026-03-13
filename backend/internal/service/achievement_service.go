package service

import (
	"context"
	"fmt"
	"log"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type AchievementService struct {
	achievements repository.AchievementRepository
}

func NewAchievementService(achievements repository.AchievementRepository) *AchievementService {
	return &AchievementService{achievements: achievements}
}

func (s *AchievementService) GetAll(ctx context.Context) ([]domain.Achievement, error) {
	achievements, err := s.achievements.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting achievements: %w", err)
	}
	return achievements, nil
}

func (s *AchievementService) GetUserAchievements(ctx context.Context, userID int64) ([]domain.UserAchievement, error) {
	achievements, err := s.achievements.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user achievements: %w", err)
	}
	return achievements, nil
}

// CheckAndUnlock checks all achievements for a user after an event.
// Returns newly unlocked achievements.
func (s *AchievementService) CheckAndUnlock(ctx context.Context, userID int64, waterGoalML int) []domain.Achievement {
	var unlocked []domain.Achievement

	allAchievements, err := s.achievements.GetAll(ctx)
	if err != nil {
		log.Printf("achievement check error: %v", err)
		return nil
	}

	for _, a := range allAchievements {
		already, err := s.achievements.IsUnlocked(ctx, userID, a.ID)
		if err != nil || already {
			continue
		}

		var currentValue int
		switch a.ConditionType {
		case "food_logs_total":
			currentValue, _ = s.achievements.CountFoodLogs(ctx, userID)
		case "tracking_streak":
			currentValue, _ = s.achievements.GetTrackingStreak(ctx, userID)
		case "water_streak":
			currentValue, _ = s.achievements.GetWaterStreak(ctx, userID, waterGoalML)
		case "unique_recipes":
			currentValue, _ = s.achievements.CountUniqueRecipes(ctx, userID)
		case "meal_plans_generated":
			currentValue, _ = s.achievements.CountMealPlans(ctx, userID)
		case "barcodes_scanned":
			currentValue, _ = s.achievements.CountBarcodeScans(ctx, userID)
		default:
			continue
		}

		if currentValue >= a.ConditionValue {
			if err := s.achievements.Unlock(ctx, userID, a.ID); err != nil {
				log.Printf("failed to unlock achievement %s: %v", a.Name, err)
				continue
			}
			unlocked = append(unlocked, a)
		}
	}

	return unlocked
}
