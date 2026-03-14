package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type PriceService struct {
	repo repository.PriceRepository
}

func NewPriceService(repo repository.PriceRepository) *PriceService {
	return &PriceService{repo: repo}
}

func (s *PriceService) LogPrice(ctx context.Context, userID int64, req domain.PriceLogRequest) (*domain.PriceLog, error) {
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}
	if req.Currency == "" {
		req.Currency = "EUR"
	}

	log := &domain.PriceLog{
		UserID:     userID,
		ItemName:   req.ItemName,
		PriceCents: req.PriceCents,
		Currency:   req.Currency,
		Store:      req.Store,
		Date:       req.Date,
	}
	if err := s.repo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("log price: %w", err)
	}
	return log, nil
}

func (s *PriceService) GetLogs(ctx context.Context, userID int64, limit int) ([]domain.PriceLog, error) {
	return s.repo.GetByUser(ctx, userID, limit)
}

func (s *PriceService) GetTrend(ctx context.Context, userID int64, itemName string) (*domain.PriceTrend, error) {
	points, err := s.repo.GetTrend(ctx, userID, itemName)
	if err != nil {
		return nil, fmt.Errorf("get trend: %w", err)
	}
	return &domain.PriceTrend{
		ItemName: itemName,
		Points:   points,
	}, nil
}

func (s *PriceService) CompareStores(ctx context.Context, userID int64, itemName string) (*domain.StoreComparison, error) {
	stores, err := s.repo.GetStoreComparison(ctx, userID, itemName)
	if err != nil {
		return nil, fmt.Errorf("compare stores: %w", err)
	}
	return &domain.StoreComparison{
		ItemName: itemName,
		Stores:   stores,
	}, nil
}

func (s *PriceService) GetSpendingSummary(ctx context.Context, userID int64, month string) (*domain.SpendingSummary, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	return s.repo.GetSpendingSummary(ctx, userID, month)
}
