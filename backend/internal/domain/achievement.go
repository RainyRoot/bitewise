package domain

import "time"

type Achievement struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Icon           string `json:"icon"`
	Category       string `json:"category"`
	ConditionType  string `json:"condition_type"`
	ConditionValue int    `json:"condition_value"`
}

type UserAchievement struct {
	UserID        int64        `json:"user_id"`
	AchievementID int64        `json:"achievement_id"`
	Achievement   *Achievement `json:"achievement,omitempty"`
	UnlockedAt    time.Time    `json:"unlocked_at"`
}
