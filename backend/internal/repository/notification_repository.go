package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type NotificationRepository interface {
	GetByUser(ctx context.Context, userID int64) ([]domain.UserNotification, error)
	Upsert(ctx context.Context, notif *domain.UserNotification) error
}

type SQLiteNotificationRepository struct {
	db *sql.DB
}

func NewSQLiteNotificationRepository(db *sql.DB) *SQLiteNotificationRepository {
	return &SQLiteNotificationRepository{db: db}
}

func (r *SQLiteNotificationRepository) GetByUser(ctx context.Context, userID int64) ([]domain.UserNotification, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, type, time, enabled FROM user_notifications WHERE user_id = ? ORDER BY type`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying notifications: %w", err)
	}
	defer rows.Close()

	var notifs []domain.UserNotification
	for rows.Next() {
		var n domain.UserNotification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Time, &n.Enabled); err != nil {
			return nil, fmt.Errorf("scanning notification: %w", err)
		}
		notifs = append(notifs, n)
	}
	return notifs, rows.Err()
}

func (r *SQLiteNotificationRepository) Upsert(ctx context.Context, notif *domain.UserNotification) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO user_notifications (user_id, type, time, enabled)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(user_id, type) DO UPDATE SET time = excluded.time, enabled = excluded.enabled`,
		notif.UserID, notif.Type, notif.Time, notif.Enabled,
	)
	if err != nil {
		return fmt.Errorf("upserting notification: %w", err)
	}
	if notif.ID == 0 {
		id, _ := result.LastInsertId()
		notif.ID = id
	}
	return nil
}
