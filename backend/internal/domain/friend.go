package domain

import "time"

type Friendship struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	FriendID  int64     `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendInvite struct {
	ID         int64     `json:"id"`
	FromUserID int64     `json:"from_user_id"`
	ToEmail    string    `json:"to_email"`
	ToUserID   *int64    `json:"to_user_id,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type FriendInviteRequest struct {
	Email string `json:"email"`
}

type FriendInfo struct {
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AddedAt   time.Time `json:"added_at"`
}

type LeaderboardEntry struct {
	UserID        int64  `json:"user_id"`
	Name          string `json:"name"`
	WeekCalories  int    `json:"week_calories"`
	CurrentStreak int    `json:"current_streak"`
	Rank          int    `json:"rank"`
}
