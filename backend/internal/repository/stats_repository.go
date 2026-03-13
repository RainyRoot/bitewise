package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type StatsRepository interface {
	GetWeeklyStats(ctx context.Context, userID int64, startDate, endDate string) ([]domain.DailyStat, error)
	GetMonthlyStats(ctx context.Context, userID int64, startDate, endDate string) ([]domain.DailyStat, error)
	GetStreaks(ctx context.Context, userID int64) (*domain.StreakInfo, error)
}

type SQLiteStatsRepository struct {
	db *sql.DB
}

func NewSQLiteStatsRepository(db *sql.DB) *SQLiteStatsRepository {
	return &SQLiteStatsRepository{db: db}
}

func (r *SQLiteStatsRepository) GetWeeklyStats(ctx context.Context, userID int64, startDate, endDate string) ([]domain.DailyStat, error) {
	return r.getDailyStats(ctx, userID, startDate, endDate)
}

func (r *SQLiteStatsRepository) GetMonthlyStats(ctx context.Context, userID int64, startDate, endDate string) ([]domain.DailyStat, error) {
	return r.getDailyStats(ctx, userID, startDate, endDate)
}

func (r *SQLiteStatsRepository) getDailyStats(ctx context.Context, userID int64, startDate, endDate string) ([]domain.DailyStat, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT f.date,
		        COALESCE(SUM(f.calories), 0),
		        COALESCE(SUM(f.protein_g), 0),
		        COALESCE(SUM(f.carbs_g), 0),
		        COALESCE(SUM(f.fat_g), 0),
		        COALESCE(SUM(f.fiber_g), 0),
		        COALESCE((SELECT SUM(w.amount_ml) FROM water_logs w WHERE w.user_id = f.user_id AND w.date = f.date), 0)
		 FROM food_logs f
		 WHERE f.user_id = ? AND f.date >= ? AND f.date <= ?
		 GROUP BY f.date
		 ORDER BY f.date`, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("querying daily stats: %w", err)
	}
	defer rows.Close()

	var stats []domain.DailyStat
	for rows.Next() {
		var s domain.DailyStat
		if err := rows.Scan(&s.Date, &s.Calories, &s.ProteinG, &s.CarbsG, &s.FatG, &s.FiberG, &s.WaterML); err != nil {
			return nil, fmt.Errorf("scanning daily stat: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (r *SQLiteStatsRepository) GetStreaks(ctx context.Context, userID int64) (*domain.StreakInfo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT date FROM food_logs WHERE user_id = ? ORDER BY date DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying streak dates: %w", err)
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, fmt.Errorf("scanning streak date: %w", err)
		}
		dates = append(dates, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	info := &domain.StreakInfo{}
	if len(dates) == 0 {
		return info, nil
	}

	// Calculate current streak
	info.CurrentStreak = calculateCurrentStreak(dates)
	info.LongestStreak = calculateLongestStreak(dates)

	return info, nil
}

func calculateCurrentStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	streak := 1
	for i := 1; i < len(dates); i++ {
		prev := parseDate(dates[i-1])
		curr := parseDate(dates[i])
		if prev.AddDate(0, 0, -1).Format("2006-01-02") == curr.Format("2006-01-02") {
			streak++
		} else {
			break
		}
	}
	return streak
}

func calculateLongestStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	longest := 1
	current := 1
	for i := 1; i < len(dates); i++ {
		prev := parseDate(dates[i-1])
		curr := parseDate(dates[i])
		if prev.AddDate(0, 0, -1).Format("2006-01-02") == curr.Format("2006-01-02") {
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 1
		}
	}
	return longest
}
