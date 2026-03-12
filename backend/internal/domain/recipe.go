package domain

import (
	"encoding/json"
	"time"
)

type Recipe struct {
	ID                 int64           `json:"id"`
	Title              string          `json:"title"`
	Description        string          `json:"description"`
	SourceURL          string          `json:"source_url"`
	SourceSite         string          `json:"source_site"`
	PrepTimeMin        int             `json:"prep_time_min"`
	CookTimeMin        int             `json:"cook_time_min"`
	Servings           int             `json:"servings"`
	Difficulty         string          `json:"difficulty"`
	ImageURL           string          `json:"image_url"`
	CaloriesPerServing int             `json:"calories_per_serving"`
	ProteinG           float64         `json:"protein_g"`
	CarbsG             float64         `json:"carbs_g"`
	FatG               float64         `json:"fat_g"`
	FiberG             float64         `json:"fiber_g"`
	Micronutrients     json.RawMessage `json:"micronutrients,omitempty"`
	Allergens          []string        `json:"allergens,omitempty"`
	Categories         []string        `json:"categories,omitempty"`
	Ingredients        []Ingredient    `json:"ingredients,omitempty"`
	IsFavorite         bool            `json:"is_favorite,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type RecipeFilter struct {
	Query        string   `json:"q"`
	Category     string   `json:"category"`
	AllergenFree []string `json:"allergen_free"`
	MaxPrepTime  int      `json:"max_prep_time"`
	Limit        int      `json:"limit"`
	Offset       int      `json:"offset"`
}
