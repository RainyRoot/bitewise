package service

import (
	"context"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

// ProfileService handles user profile operations.
type ProfileService struct {
	users repository.UserRepository
}

// NewProfileService creates a new ProfileService.
func NewProfileService(users repository.UserRepository) *ProfileService {
	return &ProfileService{users: users}
}

// GetProfile returns the user profile for the given user ID.
func (s *ProfileService) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}
	return user, nil
}

// UpdateProfile updates the user profile and returns the updated user.
func (s *ProfileService) UpdateProfile(ctx context.Context, userID int64, req domain.ProfileUpdateRequest) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user for update: %w", err)
	}

	user.Name = req.Name
	user.HeightCM = req.HeightCM
	user.WeightKG = req.WeightKG
	user.Age = req.Age
	user.Gender = req.Gender
	user.ActivityLevel = req.ActivityLevel
	user.DailyCalorieGoal = req.DailyCalorieGoal
	user.DailyWaterMLGoal = req.DailyWaterMLGoal

	if err := s.users.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("updating profile: %w", err)
	}
	return user, nil
}

// GetAllergies returns all allergies for the given user.
func (s *ProfileService) GetAllergies(ctx context.Context, userID int64) ([]domain.UserAllergy, error) {
	allergies, err := s.users.GetAllergies(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting allergies: %w", err)
	}
	return allergies, nil
}

// SetAllergies replaces all allergies for the given user.
func (s *ProfileService) SetAllergies(ctx context.Context, userID int64, allergens []string) error {
	if err := s.users.SetAllergies(ctx, userID, allergens); err != nil {
		return fmt.Errorf("setting allergies: %w", err)
	}
	return nil
}

// GetPreferences returns all preferences for the given user.
func (s *ProfileService) GetPreferences(ctx context.Context, userID int64) ([]domain.UserPreference, error) {
	prefs, err := s.users.GetPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting preferences: %w", err)
	}
	return prefs, nil
}

// SetPreferences replaces all preferences for the given user.
func (s *ProfileService) SetPreferences(ctx context.Context, userID int64, prefs []domain.UserPreference) error {
	if err := s.users.SetPreferences(ctx, userID, prefs); err != nil {
		return fmt.Errorf("setting preferences: %w", err)
	}
	return nil
}
