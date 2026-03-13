package domain

type DailyStat struct {
	Date     string  `json:"date"`
	Calories int     `json:"calories"`
	ProteinG float64 `json:"protein_g"`
	CarbsG   float64 `json:"carbs_g"`
	FatG     float64 `json:"fat_g"`
	FiberG   float64 `json:"fiber_g"`
	WaterML  int     `json:"water_ml"`
}

type StreakInfo struct {
	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`
}

type WeeklyStats struct {
	Days     []DailyStat `json:"days"`
	AvgCals  int         `json:"avg_calories"`
	AvgWater int         `json:"avg_water_ml"`
}

type MonthlyStats struct {
	Days     []DailyStat `json:"days"`
	AvgCals  int         `json:"avg_calories"`
	AvgWater int         `json:"avg_water_ml"`
}
