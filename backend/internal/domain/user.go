package domain

import "time"

type User struct {
	ID               int64     `json:"id"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	Name             string    `json:"name"`
	HeightCM         int       `json:"height_cm"`
	WeightKG         float64   `json:"weight_kg"`
	Age              int       `json:"age"`
	Gender           string    `json:"gender"`
	ActivityLevel    string    `json:"activity_level"`
	DailyCalorieGoal int       `json:"daily_calorie_goal"`
	DailyWaterMLGoal int       `json:"daily_water_ml_goal"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UserAllergy struct {
	UserID   int64  `json:"user_id"`
	Allergen string `json:"allergen"`
}

type UserPreference struct {
	UserID         int64  `json:"user_id"`
	PreferenceType string `json:"preference_type"` // "favorite" or "dislike"
	FoodItem       string `json:"food_item"`
	Category       string `json:"category"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ProfileUpdateRequest struct {
	Name             string  `json:"name"`
	HeightCM         int     `json:"height_cm"`
	WeightKG         float64 `json:"weight_kg"`
	Age              int     `json:"age"`
	Gender           string  `json:"gender"`
	ActivityLevel    string  `json:"activity_level"`
	DailyCalorieGoal int     `json:"daily_calorie_goal"`
	DailyWaterMLGoal int     `json:"daily_water_ml_goal"`
}
