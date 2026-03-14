package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type DiaryService struct {
	repo repository.DiaryRepository
}

func NewDiaryService(repo repository.DiaryRepository) *DiaryService {
	return &DiaryService{repo: repo}
}

func (s *DiaryService) CreateOrUpdate(ctx context.Context, userID int64, req domain.DiaryEntryRequest) (*domain.DiaryEntry, error) {
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	existing, err := s.repo.GetByDate(ctx, userID, req.Date)
	if err != nil {
		return nil, fmt.Errorf("check existing: %w", err)
	}

	if existing != nil {
		existing.Mood = req.Mood
		existing.EnergyLevel = req.EnergyLevel
		existing.Notes = req.Notes
		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("update diary: %w", err)
		}
		return existing, nil
	}

	entry := &domain.DiaryEntry{
		UserID:      userID,
		Date:        req.Date,
		Mood:        req.Mood,
		EnergyLevel: req.EnergyLevel,
		Notes:       req.Notes,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("create diary: %w", err)
	}
	return entry, nil
}

func (s *DiaryService) GetByDate(ctx context.Context, userID int64, date string) (*domain.DiaryEntry, error) {
	return s.repo.GetByDate(ctx, userID, date)
}

func (s *DiaryService) GetMonthly(ctx context.Context, userID int64, month string) (*domain.MonthlyMoodSummary, error) {
	from := month + "-01"
	// Parse to find last day of month
	t, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, fmt.Errorf("parse month: %w", err)
	}
	to := t.AddDate(0, 1, -1).Format("2006-01-02")

	entries, err := s.repo.GetByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get entries: %w", err)
	}

	summary := &domain.MonthlyMoodSummary{
		Month:      month,
		Entries:    entries,
		MoodCounts: make(map[string]int),
	}

	var totalEnergy float64
	for _, e := range entries {
		summary.MoodCounts[e.Mood]++
		totalEnergy += float64(e.EnergyLevel)
	}
	if len(entries) > 0 {
		summary.AvgEnergy = totalEnergy / float64(len(entries))
	}

	return summary, nil
}

func (s *DiaryService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}
