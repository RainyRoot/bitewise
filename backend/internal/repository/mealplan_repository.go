package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type MealPlanRepository interface {
	Create(ctx context.Context, plan *domain.MealPlan) error
	GetByID(ctx context.Context, id int64) (*domain.MealPlan, error)
	GetCurrent(ctx context.Context, userID int64) (*domain.MealPlan, error)
	CreateEntries(ctx context.Context, entries []domain.MealPlanEntry) error
	GetEntries(ctx context.Context, planID int64) ([]domain.MealPlanEntry, error)
	UpdateEntry(ctx context.Context, entry *domain.MealPlanEntry) error
	SetEntryLock(ctx context.Context, entryID int64, locked bool) error
	DeleteUnlockedEntries(ctx context.Context, planID int64) error
	SetStatus(ctx context.Context, planID int64, status string) error
}

type SQLiteMealPlanRepository struct {
	db *sql.DB
}

func NewSQLiteMealPlanRepository(db *sql.DB) *SQLiteMealPlanRepository {
	return &SQLiteMealPlanRepository{db: db}
}

func (r *SQLiteMealPlanRepository) Create(ctx context.Context, plan *domain.MealPlan) error {
	plan.CreatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO meal_plans (user_id, week_start_date, status, created_at) VALUES (?, ?, ?, ?)`,
		plan.UserID, plan.WeekStartDate, plan.Status, plan.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting meal plan: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	plan.ID = id
	return nil
}

func (r *SQLiteMealPlanRepository) GetByID(ctx context.Context, id int64) (*domain.MealPlan, error) {
	var plan domain.MealPlan
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, week_start_date, status, created_at FROM meal_plans WHERE id = ?`, id,
	).Scan(&plan.ID, &plan.UserID, &plan.WeekStartDate, &plan.Status, &plan.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("meal plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying meal plan: %w", err)
	}

	entries, err := r.GetEntries(ctx, plan.ID)
	if err != nil {
		return nil, err
	}
	plan.Entries = entries

	return &plan, nil
}

func (r *SQLiteMealPlanRepository) GetCurrent(ctx context.Context, userID int64) (*domain.MealPlan, error) {
	var plan domain.MealPlan
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, week_start_date, status, created_at FROM meal_plans
		 WHERE user_id = ? AND status IN ('draft', 'active')
		 ORDER BY week_start_date DESC LIMIT 1`, userID,
	).Scan(&plan.ID, &plan.UserID, &plan.WeekStartDate, &plan.Status, &plan.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no current meal plan")
	}
	if err != nil {
		return nil, fmt.Errorf("querying current meal plan: %w", err)
	}

	entries, err := r.GetEntries(ctx, plan.ID)
	if err != nil {
		return nil, err
	}
	plan.Entries = entries

	return &plan, nil
}

func (r *SQLiteMealPlanRepository) CreateEntries(ctx context.Context, entries []domain.MealPlanEntry) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	for i := range entries {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO meal_plan_entries (meal_plan_id, day_of_week, meal_type, recipe_id, servings, is_locked) VALUES (?, ?, ?, ?, ?, ?)`,
			entries[i].MealPlanID, entries[i].DayOfWeek, entries[i].MealType,
			entries[i].RecipeID, entries[i].Servings, boolToInt(entries[i].IsLocked),
		)
		if err != nil {
			return fmt.Errorf("inserting meal plan entry: %w", err)
		}
		id, _ := result.LastInsertId()
		entries[i].ID = id
	}

	return tx.Commit()
}

func (r *SQLiteMealPlanRepository) GetEntries(ctx context.Context, planID int64) ([]domain.MealPlanEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT e.id, e.meal_plan_id, e.day_of_week, e.meal_type, e.recipe_id, e.servings, e.is_locked,
		        r.id, r.title, r.description, r.prep_time_min, r.cook_time_min, r.servings, r.difficulty, r.image_url, r.calories_per_serving, r.protein_g, r.carbs_g, r.fat_g, r.fiber_g, r.allergens, r.categories, r.ingredients
		 FROM meal_plan_entries e
		 JOIN recipes r ON r.id = e.recipe_id
		 WHERE e.meal_plan_id = ?
		 ORDER BY e.day_of_week, e.meal_type`, planID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying entries: %w", err)
	}
	defer rows.Close()

	var entries []domain.MealPlanEntry
	for rows.Next() {
		var e domain.MealPlanEntry
		var rec domain.Recipe
		var locked int
		var allergensStr, categoriesStr, ingredientsStr sql.NullString

		if err := rows.Scan(
			&e.ID, &e.MealPlanID, &e.DayOfWeek, &e.MealType, &e.RecipeID, &e.Servings, &locked,
			&rec.ID, &rec.Title, &rec.Description, &rec.PrepTimeMin, &rec.CookTimeMin,
			&rec.Servings, &rec.Difficulty, &rec.ImageURL, &rec.CaloriesPerServing,
			&rec.ProteinG, &rec.CarbsG, &rec.FatG, &rec.FiberG,
			&allergensStr, &categoriesStr, &ingredientsStr,
		); err != nil {
			return nil, fmt.Errorf("scanning entry: %w", err)
		}

		e.IsLocked = locked != 0
		if allergensStr.Valid {
			json.Unmarshal([]byte(allergensStr.String), &rec.Allergens)
		}
		if categoriesStr.Valid {
			json.Unmarshal([]byte(categoriesStr.String), &rec.Categories)
		}
		if ingredientsStr.Valid {
			json.Unmarshal([]byte(ingredientsStr.String), &rec.Ingredients)
		}
		e.Recipe = &rec
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

func (r *SQLiteMealPlanRepository) UpdateEntry(ctx context.Context, entry *domain.MealPlanEntry) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE meal_plan_entries SET recipe_id = ?, servings = ? WHERE id = ?`,
		entry.RecipeID, entry.Servings, entry.ID,
	)
	if err != nil {
		return fmt.Errorf("updating entry: %w", err)
	}
	return nil
}

func (r *SQLiteMealPlanRepository) SetEntryLock(ctx context.Context, entryID int64, locked bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE meal_plan_entries SET is_locked = ? WHERE id = ?`,
		boolToInt(locked), entryID,
	)
	if err != nil {
		return fmt.Errorf("setting entry lock: %w", err)
	}
	return nil
}

func (r *SQLiteMealPlanRepository) DeleteUnlockedEntries(ctx context.Context, planID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM meal_plan_entries WHERE meal_plan_id = ? AND is_locked = 0`, planID,
	)
	if err != nil {
		return fmt.Errorf("deleting unlocked entries: %w", err)
	}
	return nil
}

func (r *SQLiteMealPlanRepository) SetStatus(ctx context.Context, planID int64, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE meal_plans SET status = ? WHERE id = ?`, status, planID,
	)
	if err != nil {
		return fmt.Errorf("setting status: %w", err)
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
