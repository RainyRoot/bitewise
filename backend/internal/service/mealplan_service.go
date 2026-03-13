package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type MealPlanService struct {
	plans   repository.MealPlanRepository
	recipes repository.RecipeRepository
	users   repository.UserRepository
}

func NewMealPlanService(plans repository.MealPlanRepository, recipes repository.RecipeRepository, users repository.UserRepository) *MealPlanService {
	return &MealPlanService{plans: plans, recipes: recipes, users: users}
}

func (s *MealPlanService) Generate(ctx context.Context, userID int64, req domain.MealPlanGenerateRequest) (*domain.MealPlan, error) {
	weekStart := req.WeekStartDate
	if weekStart == "" {
		// Default to next Monday
		now := time.Now()
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		if daysUntilMonday == 0 {
			daysUntilMonday = 7
		}
		weekStart = now.AddDate(0, 0, daysUntilMonday).Format("2006-01-02")
	}

	// Get user allergies
	allergies, err := s.users.GetAllergies(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting allergies: %w", err)
	}
	allergenList := make([]string, len(allergies))
	for i, a := range allergies {
		allergenList[i] = a.Allergen
	}

	// Get recipes excluding allergens (need 21 for 7 days × 3 meals)
	recipes, err := s.recipes.GetRandomExcludingAllergens(ctx, allergenList, 21)
	if err != nil {
		return nil, fmt.Errorf("getting recipes: %w", err)
	}
	if len(recipes) == 0 {
		return nil, fmt.Errorf("no suitable recipes found")
	}

	plan := &domain.MealPlan{
		UserID:        userID,
		WeekStartDate: weekStart,
		Status:        "draft",
	}

	if err := s.plans.Create(ctx, plan); err != nil {
		return nil, fmt.Errorf("creating meal plan: %w", err)
	}

	mealTypes := []string{"breakfast", "lunch", "dinner"}
	var entries []domain.MealPlanEntry
	idx := 0

	for day := 0; day < 7; day++ {
		for _, mealType := range mealTypes {
			recipe := recipes[idx%len(recipes)]
			entries = append(entries, domain.MealPlanEntry{
				MealPlanID: plan.ID,
				DayOfWeek:  day,
				MealType:   mealType,
				RecipeID:   recipe.ID,
				Servings:   1,
			})
			idx++
		}
	}

	if err := s.plans.CreateEntries(ctx, entries); err != nil {
		return nil, fmt.Errorf("creating entries: %w", err)
	}

	return s.plans.GetByID(ctx, plan.ID)
}

func (s *MealPlanService) GetCurrent(ctx context.Context, userID int64) (*domain.MealPlan, error) {
	plan, err := s.plans.GetCurrent(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting current plan: %w", err)
	}
	return plan, nil
}

func (s *MealPlanService) GetByID(ctx context.Context, userID, planID int64) (*domain.MealPlan, error) {
	plan, err := s.plans.GetByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("getting plan: %w", err)
	}
	if plan.UserID != userID {
		return nil, fmt.Errorf("meal plan not found")
	}
	return plan, nil
}

func (s *MealPlanService) UpdateEntry(ctx context.Context, userID, planID, entryID int64, req domain.MealPlanEntryUpdateRequest) error {
	plan, err := s.plans.GetByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("getting plan: %w", err)
	}
	if plan.UserID != userID {
		return fmt.Errorf("meal plan not found")
	}

	var entry *domain.MealPlanEntry
	for i := range plan.Entries {
		if plan.Entries[i].ID == entryID {
			entry = &plan.Entries[i]
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("entry not found")
	}

	if req.RecipeID != nil {
		entry.RecipeID = *req.RecipeID
	}
	if req.Servings != nil {
		entry.Servings = *req.Servings
	}

	return s.plans.UpdateEntry(ctx, entry)
}

func (s *MealPlanService) Activate(ctx context.Context, userID, planID int64) error {
	plan, err := s.plans.GetByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("getting plan: %w", err)
	}
	if plan.UserID != userID {
		return fmt.Errorf("meal plan not found")
	}

	return s.plans.SetStatus(ctx, planID, "active")
}
