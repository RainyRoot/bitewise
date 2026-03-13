package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type AchievementRepository interface {
	GetAll(ctx context.Context) ([]domain.Achievement, error)
	GetUserAchievements(ctx context.Context, userID int64) ([]domain.UserAchievement, error)
	IsUnlocked(ctx context.Context, userID, achievementID int64) (bool, error)
	Unlock(ctx context.Context, userID, achievementID int64) error
	GetByCondition(ctx context.Context, conditionType string, conditionValue int) (*domain.Achievement, error)

	// Stats queries for achievement checking
	CountFoodLogs(ctx context.Context, userID int64) (int, error)
	GetTrackingStreak(ctx context.Context, userID int64) (int, error)
	GetWaterStreak(ctx context.Context, userID int64, dailyGoalML int) (int, error)
	CountUniqueRecipes(ctx context.Context, userID int64) (int, error)
	CountMealPlans(ctx context.Context, userID int64) (int, error)
	CountBarcodeScans(ctx context.Context, userID int64) (int, error)
}

type SQLiteAchievementRepository struct {
	db *sql.DB
}

func NewSQLiteAchievementRepository(db *sql.DB) *SQLiteAchievementRepository {
	return &SQLiteAchievementRepository{db: db}
}

func (r *SQLiteAchievementRepository) GetAll(ctx context.Context) ([]domain.Achievement, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, icon, category, condition_type, condition_value FROM achievements ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("querying achievements: %w", err)
	}
	defer rows.Close()

	var achievements []domain.Achievement
	for rows.Next() {
		var a domain.Achievement
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Icon, &a.Category, &a.ConditionType, &a.ConditionValue); err != nil {
			return nil, fmt.Errorf("scanning achievement: %w", err)
		}
		achievements = append(achievements, a)
	}
	return achievements, rows.Err()
}

func (r *SQLiteAchievementRepository) GetUserAchievements(ctx context.Context, userID int64) ([]domain.UserAchievement, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT ua.user_id, ua.achievement_id, ua.unlocked_at,
		        a.id, a.name, a.description, a.icon, a.category, a.condition_type, a.condition_value
		 FROM user_achievements ua
		 JOIN achievements a ON a.id = ua.achievement_id
		 WHERE ua.user_id = ?
		 ORDER BY ua.unlocked_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying user achievements: %w", err)
	}
	defer rows.Close()

	var results []domain.UserAchievement
	for rows.Next() {
		var ua domain.UserAchievement
		var a domain.Achievement
		if err := rows.Scan(
			&ua.UserID, &ua.AchievementID, &ua.UnlockedAt,
			&a.ID, &a.Name, &a.Description, &a.Icon, &a.Category, &a.ConditionType, &a.ConditionValue,
		); err != nil {
			return nil, fmt.Errorf("scanning user achievement: %w", err)
		}
		ua.Achievement = &a
		results = append(results, ua)
	}
	return results, rows.Err()
}

func (r *SQLiteAchievementRepository) IsUnlocked(ctx context.Context, userID, achievementID int64) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_achievements WHERE user_id = ? AND achievement_id = ?`,
		userID, achievementID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking achievement unlock: %w", err)
	}
	return count > 0, nil
}

func (r *SQLiteAchievementRepository) Unlock(ctx context.Context, userID, achievementID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO user_achievements (user_id, achievement_id, unlocked_at) VALUES (?, ?, ?)`,
		userID, achievementID, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("unlocking achievement: %w", err)
	}
	return nil
}

func (r *SQLiteAchievementRepository) GetByCondition(ctx context.Context, conditionType string, conditionValue int) (*domain.Achievement, error) {
	var a domain.Achievement
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, icon, category, condition_type, condition_value
		 FROM achievements WHERE condition_type = ? AND condition_value = ?`,
		conditionType, conditionValue,
	).Scan(&a.ID, &a.Name, &a.Description, &a.Icon, &a.Category, &a.ConditionType, &a.ConditionValue)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying achievement by condition: %w", err)
	}
	return &a, nil
}

func (r *SQLiteAchievementRepository) CountFoodLogs(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM food_logs WHERE user_id = ?`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting food logs: %w", err)
	}
	return count, nil
}

func (r *SQLiteAchievementRepository) GetTrackingStreak(ctx context.Context, userID int64) (int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT date FROM food_logs WHERE user_id = ? ORDER BY date DESC`, userID,
	)
	if err != nil {
		return 0, fmt.Errorf("querying tracking dates: %w", err)
	}
	defer rows.Close()

	streak := 0
	expected := time.Now().UTC().Format("2006-01-02")
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return 0, fmt.Errorf("scanning date: %w", err)
		}
		if date == expected {
			streak++
			t, _ := time.Parse("2006-01-02", expected)
			expected = t.AddDate(0, 0, -1).Format("2006-01-02")
		} else if date < expected {
			break
		}
	}
	return streak, rows.Err()
}

func (r *SQLiteAchievementRepository) GetWaterStreak(ctx context.Context, userID int64, dailyGoalML int) (int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT date, SUM(amount_ml) as total FROM water_logs WHERE user_id = ? GROUP BY date ORDER BY date DESC`, userID,
	)
	if err != nil {
		return 0, fmt.Errorf("querying water dates: %w", err)
	}
	defer rows.Close()

	streak := 0
	expected := time.Now().UTC().Format("2006-01-02")
	for rows.Next() {
		var date string
		var total int
		if err := rows.Scan(&date, &total); err != nil {
			return 0, fmt.Errorf("scanning water date: %w", err)
		}
		if date == expected && total >= dailyGoalML {
			streak++
			t, _ := time.Parse("2006-01-02", expected)
			expected = t.AddDate(0, 0, -1).Format("2006-01-02")
		} else if date < expected || total < dailyGoalML {
			break
		}
	}
	return streak, rows.Err()
}

func (r *SQLiteAchievementRepository) CountUniqueRecipes(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT food_name) FROM food_logs WHERE user_id = ?`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting unique recipes: %w", err)
	}
	return count, nil
}

func (r *SQLiteAchievementRepository) CountMealPlans(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM meal_plans WHERE user_id = ?`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting meal plans: %w", err)
	}
	return count, nil
}

func (r *SQLiteAchievementRepository) CountBarcodeScans(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM food_logs WHERE user_id = ? AND barcode != ''`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting barcode scans: %w", err)
	}
	return count, nil
}
