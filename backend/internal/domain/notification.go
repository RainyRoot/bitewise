package domain

type UserNotification struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Type    string `json:"type"`
	Time    string `json:"time"`
	Enabled bool   `json:"enabled"`
}

type NotificationSettings struct {
	Notifications []UserNotification `json:"notifications"`
}
