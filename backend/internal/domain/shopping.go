package domain

import "time"

type ShoppingList struct {
	ID         int64              `json:"id"`
	UserID     int64              `json:"user_id"`
	MealPlanID int64              `json:"meal_plan_id"`
	Items      []ShoppingListItem `json:"items,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

type ShoppingListItem struct {
	ID             int64   `json:"id"`
	ShoppingListID int64   `json:"shopping_list_id"`
	IngredientName string  `json:"ingredient_name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Category       string  `json:"category"`
	IsChecked      bool    `json:"is_checked"`
}
