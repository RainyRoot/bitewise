package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type SharingService struct {
	sharing repository.SharingRepository
	recipes repository.RecipeRepository
}

func NewSharingService(sharing repository.SharingRepository, recipes repository.RecipeRepository) *SharingService {
	return &SharingService{sharing: sharing, recipes: recipes}
}

func (s *SharingService) ShareRecipe(ctx context.Context, recipeID, userID int64) (*domain.SharedRecipe, error) {
	code := generateShareCode()
	shared, err := s.sharing.CreateShare(ctx, recipeID, userID, code)
	if err != nil {
		return nil, fmt.Errorf("sharing recipe: %w", err)
	}
	return shared, nil
}

func (s *SharingService) GetSharedRecipe(ctx context.Context, code string) (*domain.Recipe, error) {
	shared, err := s.sharing.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("finding shared recipe: %w", err)
	}

	recipe, err := s.recipes.GetByID(ctx, shared.RecipeID)
	if err != nil {
		return nil, fmt.Errorf("getting shared recipe: %w", err)
	}
	return recipe, nil
}

func (s *SharingService) SaveSharedRecipe(ctx context.Context, code string, userID int64) (*domain.Recipe, error) {
	shared, err := s.sharing.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("finding shared recipe: %w", err)
	}

	original, err := s.recipes.GetByID(ctx, shared.RecipeID)
	if err != nil {
		return nil, fmt.Errorf("getting original recipe: %w", err)
	}

	// Create a copy for the user
	copy := *original
	copy.ID = 0
	copy.UserID = &userID
	copy.SourceURL = ""
	copy.SourceSite = "shared"

	if err := s.recipes.Create(ctx, &copy); err != nil {
		return nil, fmt.Errorf("saving shared recipe: %w", err)
	}
	return &copy, nil
}

func generateShareCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}
