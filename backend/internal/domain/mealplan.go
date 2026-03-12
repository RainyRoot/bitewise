package domain

import "time"

type MealPlan struct {
	ID            int64           `json:"id"`
	UserID        int64           `json:"user_id"`
	WeekStartDate string          `json:"week_start_date"` // YYYY-MM-DD (Monday)
	Status        string          `json:"status"`           // draft, active, completed
	Entries       []MealPlanEntry `json:"entries,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type MealPlanEntry struct {
	ID         int64   `json:"id"`
	MealPlanID int64   `json:"meal_plan_id"`
	DayOfWeek  int     `json:"day_of_week"` // 0=Monday, 6=Sunday
	MealType   string  `json:"meal_type"`   // breakfast, lunch, dinner, snack
	RecipeID   int64   `json:"recipe_id"`
	Servings   int     `json:"servings"`
	IsLocked   bool    `json:"is_locked"`
	Recipe     *Recipe `json:"recipe,omitempty"`
}

type MealPlanGenerateRequest struct {
	WeekStartDate string `json:"week_start_date"`
}

type MealPlanEntryUpdateRequest struct {
	RecipeID *int64 `json:"recipe_id,omitempty"`
	Servings *int   `json:"servings,omitempty"`
}
