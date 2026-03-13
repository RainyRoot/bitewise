package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type RecipeRepository interface {
	Create(ctx context.Context, recipe *domain.Recipe) error
	GetByID(ctx context.Context, id int64) (*domain.Recipe, error)
	Search(ctx context.Context, filter domain.RecipeFilter) ([]domain.Recipe, int, error)
	AddFavorite(ctx context.Context, userID, recipeID int64) error
	RemoveFavorite(ctx context.Context, userID, recipeID int64) error
	GetFavorites(ctx context.Context, userID int64, limit, offset int) ([]domain.Recipe, error)
	IsFavorite(ctx context.Context, userID, recipeID int64) (bool, error)
	GetRandomExcludingAllergens(ctx context.Context, allergens []string, limit int) ([]domain.Recipe, error)
	Count(ctx context.Context) (int, error)
}

type SQLiteRecipeRepository struct {
	db *sql.DB
}

func NewSQLiteRecipeRepository(db *sql.DB) *SQLiteRecipeRepository {
	return &SQLiteRecipeRepository{db: db}
}

func (r *SQLiteRecipeRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
	allergensJSON, _ := json.Marshal(recipe.Allergens)
	categoriesJSON, _ := json.Marshal(recipe.Categories)
	ingredientsJSON, _ := json.Marshal(recipe.Ingredients)

	recipe.CreatedAt = time.Now().UTC()

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO recipes (title, description, source_url, source_site, prep_time_min, cook_time_min, servings, difficulty, image_url, calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, micronutrients, allergens, categories, ingredients, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		recipe.Title, recipe.Description, recipe.SourceURL, recipe.SourceSite,
		recipe.PrepTimeMin, recipe.CookTimeMin, recipe.Servings, recipe.Difficulty,
		recipe.ImageURL, recipe.CaloriesPerServing, recipe.ProteinG, recipe.CarbsG,
		recipe.FatG, recipe.FiberG, recipe.Micronutrients,
		string(allergensJSON), string(categoriesJSON), string(ingredientsJSON),
		recipe.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting recipe: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	recipe.ID = id
	return nil
}

func (r *SQLiteRecipeRepository) GetByID(ctx context.Context, id int64) (*domain.Recipe, error) {
	var rec domain.Recipe
	var allergensStr, categoriesStr, ingredientsStr sql.NullString
	var microStr sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, source_url, source_site, prep_time_min, cook_time_min, servings, difficulty, image_url, calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, micronutrients, allergens, categories, ingredients, created_at
		 FROM recipes WHERE id = ?`, id,
	).Scan(
		&rec.ID, &rec.Title, &rec.Description, &rec.SourceURL, &rec.SourceSite,
		&rec.PrepTimeMin, &rec.CookTimeMin, &rec.Servings, &rec.Difficulty,
		&rec.ImageURL, &rec.CaloriesPerServing, &rec.ProteinG, &rec.CarbsG,
		&rec.FatG, &rec.FiberG, &microStr,
		&allergensStr, &categoriesStr, &ingredientsStr, &rec.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recipe not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying recipe: %w", err)
	}

	if microStr.Valid {
		rec.Micronutrients = json.RawMessage(microStr.String)
	}
	if allergensStr.Valid {
		json.Unmarshal([]byte(allergensStr.String), &rec.Allergens)
	}
	if categoriesStr.Valid {
		json.Unmarshal([]byte(categoriesStr.String), &rec.Categories)
	}
	if ingredientsStr.Valid {
		json.Unmarshal([]byte(ingredientsStr.String), &rec.Ingredients)
	}

	return &rec, nil
}

func (r *SQLiteRecipeRepository) Search(ctx context.Context, filter domain.RecipeFilter) ([]domain.Recipe, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	where := []string{"1=1"}
	args := []any{}

	if filter.Query != "" {
		where = append(where, "(title LIKE ? OR description LIKE ?)")
		q := "%" + filter.Query + "%"
		args = append(args, q, q)
	}

	if filter.Category != "" {
		where = append(where, "categories LIKE ?")
		args = append(args, "%\""+filter.Category+"\"%")
	}

	if filter.MaxPrepTime > 0 {
		where = append(where, "prep_time_min <= ?")
		args = append(args, filter.MaxPrepTime)
	}

	for _, allergen := range filter.AllergenFree {
		where = append(where, "allergens NOT LIKE ?")
		args = append(args, "%\""+allergen+"\"%")
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM recipes WHERE "+whereClause, args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting recipes: %w", err)
	}

	queryArgs := append(args, filter.Limit, filter.Offset)
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, title, description, source_url, source_site, prep_time_min, cook_time_min, servings, difficulty, image_url, calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, allergens, categories, ingredients, created_at FROM recipes WHERE "+whereClause+" ORDER BY id DESC LIMIT ? OFFSET ?",
		queryArgs...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("searching recipes: %w", err)
	}
	defer rows.Close()

	var recipes []domain.Recipe
	for rows.Next() {
		var rec domain.Recipe
		var allergensStr, categoriesStr, ingredientsStr sql.NullString
		if err := rows.Scan(
			&rec.ID, &rec.Title, &rec.Description, &rec.SourceURL, &rec.SourceSite,
			&rec.PrepTimeMin, &rec.CookTimeMin, &rec.Servings, &rec.Difficulty,
			&rec.ImageURL, &rec.CaloriesPerServing, &rec.ProteinG, &rec.CarbsG,
			&rec.FatG, &rec.FiberG,
			&allergensStr, &categoriesStr, &ingredientsStr, &rec.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning recipe: %w", err)
		}
		if allergensStr.Valid {
			json.Unmarshal([]byte(allergensStr.String), &rec.Allergens)
		}
		if categoriesStr.Valid {
			json.Unmarshal([]byte(categoriesStr.String), &rec.Categories)
		}
		if ingredientsStr.Valid {
			json.Unmarshal([]byte(ingredientsStr.String), &rec.Ingredients)
		}
		recipes = append(recipes, rec)
	}

	return recipes, total, rows.Err()
}

func (r *SQLiteRecipeRepository) AddFavorite(ctx context.Context, userID, recipeID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO user_favorite_recipes (user_id, recipe_id) VALUES (?, ?)`,
		userID, recipeID,
	)
	if err != nil {
		return fmt.Errorf("adding favorite: %w", err)
	}
	return nil
}

func (r *SQLiteRecipeRepository) RemoveFavorite(ctx context.Context, userID, recipeID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM user_favorite_recipes WHERE user_id = ? AND recipe_id = ?`,
		userID, recipeID,
	)
	if err != nil {
		return fmt.Errorf("removing favorite: %w", err)
	}
	return nil
}

func (r *SQLiteRecipeRepository) GetFavorites(ctx context.Context, userID int64, limit, offset int) ([]domain.Recipe, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT r.id, r.title, r.description, r.source_url, r.source_site, r.prep_time_min, r.cook_time_min, r.servings, r.difficulty, r.image_url, r.calories_per_serving, r.protein_g, r.carbs_g, r.fat_g, r.fiber_g, r.allergens, r.categories, r.ingredients, r.created_at
		 FROM recipes r
		 JOIN user_favorite_recipes f ON f.recipe_id = r.id
		 WHERE f.user_id = ?
		 ORDER BY f.created_at DESC
		 LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying favorites: %w", err)
	}
	defer rows.Close()

	var recipes []domain.Recipe
	for rows.Next() {
		var rec domain.Recipe
		var allergensStr, categoriesStr, ingredientsStr sql.NullString
		if err := rows.Scan(
			&rec.ID, &rec.Title, &rec.Description, &rec.SourceURL, &rec.SourceSite,
			&rec.PrepTimeMin, &rec.CookTimeMin, &rec.Servings, &rec.Difficulty,
			&rec.ImageURL, &rec.CaloriesPerServing, &rec.ProteinG, &rec.CarbsG,
			&rec.FatG, &rec.FiberG,
			&allergensStr, &categoriesStr, &ingredientsStr, &rec.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning favorite recipe: %w", err)
		}
		if allergensStr.Valid {
			json.Unmarshal([]byte(allergensStr.String), &rec.Allergens)
		}
		if categoriesStr.Valid {
			json.Unmarshal([]byte(categoriesStr.String), &rec.Categories)
		}
		if ingredientsStr.Valid {
			json.Unmarshal([]byte(ingredientsStr.String), &rec.Ingredients)
		}
		rec.IsFavorite = true
		recipes = append(recipes, rec)
	}

	return recipes, rows.Err()
}

func (r *SQLiteRecipeRepository) IsFavorite(ctx context.Context, userID, recipeID int64) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_favorite_recipes WHERE user_id = ? AND recipe_id = ?`,
		userID, recipeID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking favorite: %w", err)
	}
	return count > 0, nil
}

func (r *SQLiteRecipeRepository) GetRandomExcludingAllergens(ctx context.Context, allergens []string, limit int) ([]domain.Recipe, error) {
	where := "1=1"
	args := []any{}

	for _, a := range allergens {
		where += " AND allergens NOT LIKE ?"
		args = append(args, "%\""+a+"\"%")
	}

	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, title, description, source_url, source_site, prep_time_min, cook_time_min, servings, difficulty, image_url, calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, allergens, categories, ingredients, created_at FROM recipes WHERE "+where+" ORDER BY RANDOM() LIMIT ?",
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("querying random recipes: %w", err)
	}
	defer rows.Close()

	var recipes []domain.Recipe
	for rows.Next() {
		var rec domain.Recipe
		var allergensStr, categoriesStr, ingredientsStr sql.NullString
		if err := rows.Scan(
			&rec.ID, &rec.Title, &rec.Description, &rec.SourceURL, &rec.SourceSite,
			&rec.PrepTimeMin, &rec.CookTimeMin, &rec.Servings, &rec.Difficulty,
			&rec.ImageURL, &rec.CaloriesPerServing, &rec.ProteinG, &rec.CarbsG,
			&rec.FatG, &rec.FiberG,
			&allergensStr, &categoriesStr, &ingredientsStr, &rec.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning random recipe: %w", err)
		}
		if allergensStr.Valid {
			json.Unmarshal([]byte(allergensStr.String), &rec.Allergens)
		}
		if categoriesStr.Valid {
			json.Unmarshal([]byte(categoriesStr.String), &rec.Categories)
		}
		if ingredientsStr.Valid {
			json.Unmarshal([]byte(ingredientsStr.String), &rec.Ingredients)
		}
		recipes = append(recipes, rec)
	}

	return recipes, rows.Err()
}

func (r *SQLiteRecipeRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM recipes").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting recipes: %w", err)
	}
	return count, nil
}
