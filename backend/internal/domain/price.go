package domain

import "time"

type PriceLog struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ItemName   string    `json:"item_name"`
	PriceCents int       `json:"price_cents"`
	Currency   string    `json:"currency"`
	Store      string    `json:"store"`
	Date       string    `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
}

type PriceLogRequest struct {
	ItemName   string `json:"item_name"`
	PriceCents int    `json:"price_cents"`
	Currency   string `json:"currency"`
	Store      string `json:"store"`
	Date       string `json:"date"`
}

type PriceTrend struct {
	ItemName string          `json:"item_name"`
	Points   []PriceTrendPoint `json:"points"`
}

type PriceTrendPoint struct {
	Date       string `json:"date"`
	PriceCents int    `json:"price_cents"`
	Store      string `json:"store"`
}

type StoreComparison struct {
	ItemName string              `json:"item_name"`
	Stores   []StorePrice        `json:"stores"`
}

type StorePrice struct {
	Store         string `json:"store"`
	AvgPriceCents int    `json:"avg_price_cents"`
	LastPrice     int    `json:"last_price_cents"`
	EntryCount    int    `json:"entry_count"`
}

type SpendingSummary struct {
	Month      string `json:"month"`
	TotalCents int    `json:"total_cents"`
	ItemCount  int    `json:"item_count"`
}
