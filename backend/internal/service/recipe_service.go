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

func (s *RecipeService) CreateCustomRecipe(ctx context.Context, userID int64, req domain.CreateRecipeRequest) (*domain.Recipe, error) {
	recipe := &domain.Recipe{
		UserID:             &userID,
		Title:              req.Title,
		Description:        req.Description,
		PrepTimeMin:        req.PrepTimeMin,
		CookTimeMin:        req.CookTimeMin,
		Servings:           req.Servings,
		Difficulty:         req.Difficulty,
		ImageURL:           req.ImageURL,
		CaloriesPerServing: req.CaloriesPerServing,
		ProteinG:           req.ProteinG,
		CarbsG:             req.CarbsG,
		FatG:               req.FatG,
		FiberG:             req.FiberG,
		Allergens:          req.Allergens,
		Categories:         req.Categories,
		Ingredients:        req.Ingredients,
		Instructions:       req.Instructions,
		SourceSite:         "custom",
	}

	if err := s.recipes.Create(ctx, recipe); err != nil {
		return nil, fmt.Errorf("creating custom recipe: %w", err)
	}
	return recipe, nil
}

func (s *RecipeService) GetMyRecipes(ctx context.Context, userID int64, limit, offset int) ([]domain.Recipe, error) {
	recipes, err := s.recipes.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getting user recipes: %w", err)
	}
	return recipes, nil
}

func (s *RecipeService) DeleteRecipe(ctx context.Context, id, userID int64) error {
	if err := s.recipes.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("deleting recipe: %w", err)
	}
	return nil
}
