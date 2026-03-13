package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

// DataProvider defines the interface for recipe data sources.
type DataProvider interface {
	SearchRecipes(ctx context.Context, query string, limit int) ([]domain.Recipe, error)
	GetRecipeDetail(ctx context.Context, sourceURL string) (*domain.Recipe, error)
}

// ChefkochScraper fetches recipes from the Chefkoch.de API.
type ChefkochScraper struct {
	client *http.Client
}

// NewChefkochScraper creates a new ChefkochScraper.
func NewChefkochScraper() *ChefkochScraper {
	return &ChefkochScraper{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type chefkochSearchResponse struct {
	Results []chefkochRecipe `json:"results"`
}

type chefkochRecipe struct {
	Recipe struct {
		ID              string  `json:"id"`
		Title           string  `json:"title"`
		Subtitle        string  `json:"subtitle"`
		PreparationTime int     `json:"preparationTime"`
		Servings        int     `json:"servings"`
		Difficulty      int     `json:"difficulty"`
		Rating          float64 `json:"rating"`
		PreviewImageURL string  `json:"previewImageUrlTemplate"`
	} `json:"recipe"`
}

type chefkochDetailResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Subtitle        string  `json:"subtitle"`
	PreparationTime int     `json:"preparationTime"`
	Servings        int     `json:"servings"`
	Difficulty      int     `json:"difficulty"`
	Rating          float64 `json:"rating"`
	Instructions    string  `json:"instructions"`
	IngredientGroups []struct {
		Ingredients []struct {
			Name   string  `json:"name"`
			Amount float64 `json:"amount"`
			Unit   string  `json:"unit"`
		} `json:"ingredients"`
	} `json:"ingredientGroups"`
	PreviewImageURL string `json:"previewImageUrlTemplate"`
}

func (s *ChefkochScraper) SearchRecipes(ctx context.Context, query string, limit int) ([]domain.Recipe, error) {
	if limit <= 0 {
		limit = 10
	}

	reqURL := fmt.Sprintf("https://api.chefkoch.de/v2/search?query=%s&limit=%d&type=recipe",
		url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "BiteWise/1.0 (Nutrition App)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching search results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chefkoch API returned status %d", resp.StatusCode)
	}

	var searchResp chefkochSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}

	var recipes []domain.Recipe
	for _, r := range searchResp.Results {
		difficulty := "medium"
		switch r.Recipe.Difficulty {
		case 1:
			difficulty = "easy"
		case 3:
			difficulty = "hard"
		}

		recipes = append(recipes, domain.Recipe{
			Title:       r.Recipe.Title,
			Description: r.Recipe.Subtitle,
			SourceURL:   fmt.Sprintf("https://www.chefkoch.de/rezepte/%s", r.Recipe.ID),
			SourceSite:  "chefkoch.de",
			PrepTimeMin: r.Recipe.PreparationTime,
			Servings:    r.Recipe.Servings,
			Difficulty:  difficulty,
			ImageURL:    r.Recipe.PreviewImageURL,
		})
	}

	return recipes, nil
}

func (s *ChefkochScraper) GetRecipeDetail(ctx context.Context, sourceURL string) (*domain.Recipe, error) {
	// Extract recipe ID from URL — not implemented in Phase 1
	// This would parse the URL to get the ID and call the detail API
	return nil, fmt.Errorf("recipe detail fetching not yet implemented")
}
