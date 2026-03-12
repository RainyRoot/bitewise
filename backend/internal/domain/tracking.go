package domain

import "time"

type FoodLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Date      string    `json:"date"` // YYYY-MM-DD
	MealType  string    `json:"meal_type"`
	FoodName  string    `json:"food_name"`
	Barcode   string    `json:"barcode,omitempty"`
	Servings  float64   `json:"servings"`
	Calories  int       `json:"calories"`
	ProteinG  float64   `json:"protein_g"`
	CarbsG    float64   `json:"carbs_g"`
	FatG      float64   `json:"fat_g"`
	FiberG    float64   `json:"fiber_g"`
	CreatedAt time.Time `json:"created_at"`
}

type FoodLogRequest struct {
	Date     string  `json:"date"`
	MealType string  `json:"meal_type"`
	FoodName string  `json:"food_name"`
	Barcode  string  `json:"barcode,omitempty"`
	Servings float64 `json:"servings"`
	Calories int     `json:"calories"`
	ProteinG float64 `json:"protein_g"`
	CarbsG   float64 `json:"carbs_g"`
	FatG     float64 `json:"fat_g"`
	FiberG   float64 `json:"fiber_g"`
}

type WaterLog struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	Date     string    `json:"date"`
	AmountML int       `json:"amount_ml"`
	LoggedAt time.Time `json:"logged_at"`
}

type WaterLogRequest struct {
	Date     string `json:"date"`
	AmountML int    `json:"amount_ml"`
}

type NutritionSummary struct {
	Date      string  `json:"date"`
	Calories  int     `json:"calories"`
	ProteinG  float64 `json:"protein_g"`
	CarbsG    float64 `json:"carbs_g"`
	FatG      float64 `json:"fat_g"`
	FiberG    float64 `json:"fiber_g"`
	WaterML   int     `json:"water_ml"`
	MealCount int     `json:"meal_count"`
}
