package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	SetAllergies(ctx context.Context, userID int64, allergens []string) error
	GetAllergies(ctx context.Context, userID int64) ([]domain.UserAllergy, error)
	SetPreferences(ctx context.Context, userID int64, prefs []domain.UserPreference) error
	GetPreferences(ctx context.Context, userID int64) ([]domain.UserPreference, error)
}

// SQLiteUserRepository implements UserRepository using SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, name, height_cm, weight_kg, age, gender, activity_level, daily_calorie_goal, daily_water_ml_goal, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Email, user.PasswordHash, user.Name, user.HeightCM, user.WeightKG,
		user.Age, user.Gender, user.ActivityLevel, user.DailyCalorieGoal,
		user.DailyWaterMLGoal, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	user.ID = id
	return nil
}

func (r *SQLiteUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, height_cm, weight_kg, age, gender, activity_level, daily_calorie_goal, daily_water_ml_goal, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.HeightCM, &u.WeightKG,
		&u.Age, &u.Gender, &u.ActivityLevel, &u.DailyCalorieGoal,
		&u.DailyWaterMLGoal, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by id: %w", err)
	}
	return &u, nil
}

func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, height_cm, weight_kg, age, gender, activity_level, daily_calorie_goal, daily_water_ml_goal, created_at, updated_at
		 FROM users WHERE email = ?`, email,
	).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.HeightCM, &u.WeightKG,
		&u.Age, &u.Gender, &u.ActivityLevel, &u.DailyCalorieGoal,
		&u.DailyWaterMLGoal, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by email: %w", err)
	}
	return &u, nil
}

func (r *SQLiteUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now().UTC()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET name = ?, height_cm = ?, weight_kg = ?, age = ?, gender = ?, activity_level = ?, daily_calorie_goal = ?, daily_water_ml_goal = ?, updated_at = ?
		 WHERE id = ?`,
		user.Name, user.HeightCM, user.WeightKG, user.Age, user.Gender,
		user.ActivityLevel, user.DailyCalorieGoal, user.DailyWaterMLGoal,
		user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

func (r *SQLiteUserRepository) SetAllergies(ctx context.Context, userID int64, allergens []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_allergies WHERE user_id = ?`, userID); err != nil {
		return fmt.Errorf("deleting allergies: %w", err)
	}

	for _, allergen := range allergens {
		if _, err := tx.ExecContext(ctx, `INSERT INTO user_allergies (user_id, allergen) VALUES (?, ?)`, userID, allergen); err != nil {
			return fmt.Errorf("inserting allergy %q: %w", allergen, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing allergies: %w", err)
	}
	return nil
}

func (r *SQLiteUserRepository) GetAllergies(ctx context.Context, userID int64) ([]domain.UserAllergy, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT user_id, allergen FROM user_allergies WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying allergies: %w", err)
	}
	defer rows.Close()

	var allergies []domain.UserAllergy
	for rows.Next() {
		var a domain.UserAllergy
		if err := rows.Scan(&a.UserID, &a.Allergen); err != nil {
			return nil, fmt.Errorf("scanning allergy: %w", err)
		}
		allergies = append(allergies, a)
	}
	return allergies, rows.Err()
}

func (r *SQLiteUserRepository) SetPreferences(ctx context.Context, userID int64, prefs []domain.UserPreference) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_preferences WHERE user_id = ?`, userID); err != nil {
		return fmt.Errorf("deleting preferences: %w", err)
	}

	for _, p := range prefs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO user_preferences (user_id, preference_type, food_item, category) VALUES (?, ?, ?, ?)`,
			userID, p.PreferenceType, p.FoodItem, p.Category,
		); err != nil {
			return fmt.Errorf("inserting preference: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing preferences: %w", err)
	}
	return nil
}

func (r *SQLiteUserRepository) GetPreferences(ctx context.Context, userID int64) ([]domain.UserPreference, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT user_id, preference_type, food_item, category FROM user_preferences WHERE user_id = ?`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying preferences: %w", err)
	}
	defer rows.Close()

	var prefs []domain.UserPreference
	for rows.Next() {
		var p domain.UserPreference
		if err := rows.Scan(&p.UserID, &p.PreferenceType, &p.FoodItem, &p.Category); err != nil {
			return nil, fmt.Errorf("scanning preference: %w", err)
		}
		prefs = append(prefs, p)
	}
	return prefs, rows.Err()
}
