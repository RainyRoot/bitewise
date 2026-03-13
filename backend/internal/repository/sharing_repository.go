package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type SharingRepository interface {
	CreateShare(ctx context.Context, recipeID, userID int64, shareCode string) (*domain.SharedRecipe, error)
	GetByCode(ctx context.Context, code string) (*domain.SharedRecipe, error)
}

type SQLiteSharingRepository struct {
	db *sql.DB
}

func NewSQLiteSharingRepository(db *sql.DB) *SQLiteSharingRepository {
	return &SQLiteSharingRepository{db: db}
}

func (r *SQLiteSharingRepository) CreateShare(ctx context.Context, recipeID, userID int64, shareCode string) (*domain.SharedRecipe, error) {
	// Check if already shared
	var existing domain.SharedRecipe
	err := r.db.QueryRowContext(ctx,
		`SELECT id, recipe_id, share_code, created_by, created_at FROM shared_recipes WHERE recipe_id = ? AND created_by = ?`,
		recipeID, userID,
	).Scan(&existing.ID, &existing.RecipeID, &existing.ShareCode, &existing.CreatedBy, &existing.CreatedAt)
	if err == nil {
		return &existing, nil
	}

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO shared_recipes (recipe_id, share_code, created_by, created_at) VALUES (?, ?, ?, ?)`,
		recipeID, shareCode, userID, now,
	)
	if err != nil {
		return nil, fmt.Errorf("creating share: %w", err)
	}

	id, _ := result.LastInsertId()
	return &domain.SharedRecipe{
		ID:        id,
		RecipeID:  recipeID,
		ShareCode: shareCode,
		CreatedBy: userID,
		CreatedAt: now,
	}, nil
}

func (r *SQLiteSharingRepository) GetByCode(ctx context.Context, code string) (*domain.SharedRecipe, error) {
	var s domain.SharedRecipe
	err := r.db.QueryRowContext(ctx,
		`SELECT id, recipe_id, share_code, created_by, created_at FROM shared_recipes WHERE share_code = ?`, code,
	).Scan(&s.ID, &s.RecipeID, &s.ShareCode, &s.CreatedBy, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shared recipe not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying shared recipe: %w", err)
	}
	return &s, nil
}
