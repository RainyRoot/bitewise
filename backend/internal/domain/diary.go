package domain

import "time"

type DiaryEntry struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Date        string    `json:"date"`
	Mood        string    `json:"mood"`
	EnergyLevel int       `json:"energy_level"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DiaryEntryRequest struct {
	Date        string `json:"date"`
	Mood        string `json:"mood"`
	EnergyLevel int    `json:"energy_level"`
	Notes       string `json:"notes"`
}

type MonthlyMoodSummary struct {
	Month       string       `json:"month"`
	Entries     []DiaryEntry `json:"entries"`
	AvgEnergy   float64      `json:"avg_energy"`
	MoodCounts  map[string]int `json:"mood_counts"`
}
