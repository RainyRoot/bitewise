package service

import (
	"context"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type RecipeService struct {
	recipes repository.RecipeRepository
}

func NewRecipeService(recipes repository.RecipeRepository) *RecipeService {
	return &RecipeService{recipes: recipes}
}

func (s *RecipeService) Search(ctx context.Context, filter domain.RecipeFilter) ([]domain.Recipe, int, error) {
	recipes, total, err := s.recipes.Search(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("searching recipes: %w", err)
	}
	return recipes, total, nil
}

func (s *RecipeService) GetByID(ctx context.Context, id int64, userID int64) (*domain.Recipe, error) {
	recipe, err := s.recipes.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting recipe: %w", err)
	}

	if userID > 0 {
		fav, err := s.recipes.IsFavorite(ctx, userID, id)
		if err == nil {
			recipe.IsFavorite = fav
		}
	}

	return recipe, nil
}

func (s *RecipeService) AddFavorite(ctx context.Context, userID, recipeID int64) error {
	if err := s.recipes.AddFavorite(ctx, userID, recipeID); err != nil {
		return fmt.Errorf("adding favorite: %w", err)
	}
	return nil
}

func (s *RecipeService) RemoveFavorite(ctx context.Context, userID, recipeID int64) error {
	if err := s.recipes.RemoveFavorite(ctx, userID, recipeID); err != nil {
		return fmt.Errorf("removing favorite: %w", err)
	}
	return nil
}

func (s *RecipeService) GetFavorites(ctx context.Context, userID int64, limit, offset int) ([]domain.Recipe, error) {
	recipes, err := s.recipes.GetFavorites(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getting favorites: %w", err)
	}
	return recipes, nil
}
