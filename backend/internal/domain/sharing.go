package domain

import "time"

type SharedRecipe struct {
	ID        int64     `json:"id"`
	RecipeID  int64     `json:"recipe_id"`
	ShareCode string    `json:"share_code"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
