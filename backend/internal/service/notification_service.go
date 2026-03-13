package service

import (
	"context"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type NotificationService struct {
	notifications repository.NotificationRepository
}

func NewNotificationService(notifications repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifications: notifications}
}

func (s *NotificationService) GetSettings(ctx context.Context, userID int64) ([]domain.UserNotification, error) {
	notifs, err := s.notifications.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting notification settings: %w", err)
	}
	return notifs, nil
}

func (s *NotificationService) UpdateSettings(ctx context.Context, userID int64, settings []domain.UserNotification) ([]domain.UserNotification, error) {
	for i := range settings {
		settings[i].UserID = userID
		if err := s.notifications.Upsert(ctx, &settings[i]); err != nil {
			return nil, fmt.Errorf("updating notification %s: %w", settings[i].Type, err)
		}
	}

	return s.notifications.GetByUser(ctx, userID)
}
