package domain

import "time"

type PantryItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type PantryRequest struct {
	Items []string `json:"items"`
}

type PantryMatch struct {
	Recipe       Recipe `json:"recipe"`
	MatchedCount int    `json:"matched_count"`
	TotalCount   int    `json:"total_count"`
	MatchPercent int    `json:"match_percent"`
}
