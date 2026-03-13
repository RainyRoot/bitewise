package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type TrackingRepository interface {
	CreateFoodLog(ctx context.Context, log *domain.FoodLog) error
	GetFoodLogs(ctx context.Context, userID int64, date string) ([]domain.FoodLog, error)
	DeleteFoodLog(ctx context.Context, id, userID int64) error
	CreateWaterLog(ctx context.Context, log *domain.WaterLog) error
	GetWaterLogs(ctx context.Context, userID int64, date string) ([]domain.WaterLog, error)
	GetNutritionSummary(ctx context.Context, userID int64, date string) (*domain.NutritionSummary, error)
}

type SQLiteTrackingRepository struct {
	db *sql.DB
}

func NewSQLiteTrackingRepository(db *sql.DB) *SQLiteTrackingRepository {
	return &SQLiteTrackingRepository{db: db}
}

func (r *SQLiteTrackingRepository) CreateFoodLog(ctx context.Context, log *domain.FoodLog) error {
	log.CreatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO food_logs (user_id, date, meal_type, food_name, barcode, servings, calories, protein_g, carbs_g, fat_g, fiber_g, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		log.UserID, log.Date, log.MealType, log.FoodName, log.Barcode,
		log.Servings, log.Calories, log.ProteinG, log.CarbsG, log.FatG,
		log.FiberG, log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting food log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	log.ID = id
	return nil
}

func (r *SQLiteTrackingRepository) GetFoodLogs(ctx context.Context, userID int64, date string) ([]domain.FoodLog, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, date, meal_type, food_name, barcode, servings, calories, protein_g, carbs_g, fat_g, fiber_g, created_at
		 FROM food_logs WHERE user_id = ? AND date = ? ORDER BY created_at`,
		userID, date,
	)
	if err != nil {
		return nil, fmt.Errorf("querying food logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.FoodLog
	for rows.Next() {
		var l domain.FoodLog
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Date, &l.MealType, &l.FoodName, &l.Barcode,
			&l.Servings, &l.Calories, &l.ProteinG, &l.CarbsG, &l.FatG,
			&l.FiberG, &l.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning food log: %w", err)
		}
		logs = append(logs, l)
	}

	return logs, rows.Err()
}

func (r *SQLiteTrackingRepository) DeleteFoodLog(ctx context.Context, id, userID int64) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM food_logs WHERE id = ? AND user_id = ?`, id, userID,
	)
	if err != nil {
		return fmt.Errorf("deleting food log: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("food log not found")
	}
	return nil
}

func (r *SQLiteTrackingRepository) CreateWaterLog(ctx context.Context, log *domain.WaterLog) error {
	log.LoggedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO water_logs (user_id, date, amount_ml, logged_at) VALUES (?, ?, ?, ?)`,
		log.UserID, log.Date, log.AmountML, log.LoggedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting water log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	log.ID = id
	return nil
}

func (r *SQLiteTrackingRepository) GetWaterLogs(ctx context.Context, userID int64, date string) ([]domain.WaterLog, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, date, amount_ml, logged_at FROM water_logs WHERE user_id = ? AND date = ? ORDER BY logged_at`,
		userID, date,
	)
	if err != nil {
		return nil, fmt.Errorf("querying water logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.WaterLog
	for rows.Next() {
		var l domain.WaterLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Date, &l.AmountML, &l.LoggedAt); err != nil {
			return nil, fmt.Errorf("scanning water log: %w", err)
		}
		logs = append(logs, l)
	}

	return logs, rows.Err()
}

func (r *SQLiteTrackingRepository) GetNutritionSummary(ctx context.Context, userID int64, date string) (*domain.NutritionSummary, error) {
	summary := &domain.NutritionSummary{Date: date}

	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(calories), 0), COALESCE(SUM(protein_g), 0), COALESCE(SUM(carbs_g), 0), COALESCE(SUM(fat_g), 0), COALESCE(SUM(fiber_g), 0), COUNT(*)
		 FROM food_logs WHERE user_id = ? AND date = ?`,
		userID, date,
	).Scan(&summary.Calories, &summary.ProteinG, &summary.CarbsG, &summary.FatG, &summary.FiberG, &summary.MealCount)
	if err != nil {
		return nil, fmt.Errorf("querying food summary: %w", err)
	}

	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount_ml), 0) FROM water_logs WHERE user_id = ? AND date = ?`,
		userID, date,
	).Scan(&summary.WaterML)
	if err != nil {
		return nil, fmt.Errorf("querying water summary: %w", err)
	}

	return summary, nil
}
