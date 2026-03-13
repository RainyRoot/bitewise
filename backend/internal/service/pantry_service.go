package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type PantryService struct {
	pantry  repository.PantryRepository
	recipes repository.RecipeRepository
}

func NewPantryService(pantry repository.PantryRepository, recipes repository.RecipeRepository) *PantryService {
	return &PantryService{pantry: pantry, recipes: recipes}
}

func (s *PantryService) SetItems(ctx context.Context, userID int64, items []string) ([]domain.PantryItem, error) {
	if err := s.pantry.SetItems(ctx, userID, items); err != nil {
		return nil, fmt.Errorf("setting pantry items: %w", err)
	}
	return s.pantry.GetItems(ctx, userID)
}

func (s *PantryService) GetItems(ctx context.Context, userID int64) ([]domain.PantryItem, error) {
	return s.pantry.GetItems(ctx, userID)
}

func (s *PantryService) FindRecipes(ctx context.Context, userID int64) ([]domain.PantryMatch, error) {
	pantryItems, err := s.pantry.GetItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting pantry items: %w", err)
	}
	if len(pantryItems) == 0 {
		return nil, fmt.Errorf("pantry is empty")
	}

	pantryNames := make([]string, len(pantryItems))
	for i, item := range pantryItems {
		pantryNames[i] = strings.ToLower(item.Name)
	}

	// Get all recipes (limited to 200)
	allRecipes, _, err := s.recipes.Search(ctx, domain.RecipeFilter{Limit: 200})
	if err != nil {
		return nil, fmt.Errorf("searching recipes: %w", err)
	}

	var matches []domain.PantryMatch
	for _, recipe := range allRecipes {
		if len(recipe.Ingredients) == 0 {
			continue
		}

		matched := 0
		for _, ing := range recipe.Ingredients {
			ingLower := strings.ToLower(ing.Name)
			for _, pantryName := range pantryNames {
				if strings.Contains(ingLower, pantryName) || strings.Contains(pantryName, ingLower) {
					matched++
					break
				}
			}
		}

		if matched == 0 {
			continue
		}

		percent := (matched * 100) / len(recipe.Ingredients)
		matches = append(matches, domain.PantryMatch{
			Recipe:       recipe,
			MatchedCount: matched,
			TotalCount:   len(recipe.Ingredients),
			MatchPercent: percent,
		})
	}

	// Sort by match percentage descending
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchPercent > matches[j].MatchPercent
	})

	// Return top 20
	if len(matches) > 20 {
		matches = matches[:20]
	}

	return matches, nil
}
