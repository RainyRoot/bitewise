package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type PriceRepository interface {
	Create(ctx context.Context, log *domain.PriceLog) error
	GetByUser(ctx context.Context, userID int64, limit int) ([]domain.PriceLog, error)
	GetTrend(ctx context.Context, userID int64, itemName string) ([]domain.PriceTrendPoint, error)
	GetStoreComparison(ctx context.Context, userID int64, itemName string) ([]domain.StorePrice, error)
	GetSpendingSummary(ctx context.Context, userID int64, month string) (*domain.SpendingSummary, error)
}

type sqlitePriceRepo struct {
	db *sql.DB
}

func NewSQLitePriceRepository(db *sql.DB) PriceRepository {
	return &sqlitePriceRepo{db: db}
}

func (r *sqlitePriceRepo) Create(ctx context.Context, log *domain.PriceLog) error {
	if log.Currency == "" {
		log.Currency = "EUR"
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO price_logs (user_id, item_name, price_cents, currency, store, date)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		log.UserID, log.ItemName, log.PriceCents, log.Currency, log.Store, log.Date,
	)
	if err != nil {
		return fmt.Errorf("insert price log: %w", err)
	}
	id, _ := res.LastInsertId()
	log.ID = id
	return nil
}

func (r *sqlitePriceRepo) GetByUser(ctx context.Context, userID int64, limit int) ([]domain.PriceLog, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, item_name, price_cents, currency, store, date, created_at
		 FROM price_logs WHERE user_id = ? ORDER BY date DESC, created_at DESC LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query price logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.PriceLog
	for rows.Next() {
		var l domain.PriceLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.ItemName, &l.PriceCents, &l.Currency, &l.Store, &l.Date, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan price log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *sqlitePriceRepo) GetTrend(ctx context.Context, userID int64, itemName string) ([]domain.PriceTrendPoint, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT date, price_cents, store FROM price_logs
		 WHERE user_id = ? AND LOWER(item_name) = LOWER(?)
		 ORDER BY date ASC`,
		userID, itemName,
	)
	if err != nil {
		return nil, fmt.Errorf("query price trend: %w", err)
	}
	defer rows.Close()

	var points []domain.PriceTrendPoint
	for rows.Next() {
		var p domain.PriceTrendPoint
		if err := rows.Scan(&p.Date, &p.PriceCents, &p.Store); err != nil {
			return nil, fmt.Errorf("scan price trend: %w", err)
		}
		points = append(points, p)
	}
	return points, nil
}

func (r *sqlitePriceRepo) GetStoreComparison(ctx context.Context, userID int64, itemName string) ([]domain.StorePrice, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT store, CAST(AVG(price_cents) AS INTEGER), COUNT(*)
		 FROM price_logs
		 WHERE user_id = ? AND LOWER(item_name) = LOWER(?) AND store != ''
		 GROUP BY store
		 ORDER BY AVG(price_cents) ASC`,
		userID, itemName,
	)
	if err != nil {
		return nil, fmt.Errorf("query store comparison: %w", err)
	}
	defer rows.Close()

	var stores []domain.StorePrice
	for rows.Next() {
		var s domain.StorePrice
		if err := rows.Scan(&s.Store, &s.AvgPriceCents, &s.EntryCount); err != nil {
			return nil, fmt.Errorf("scan store comparison: %w", err)
		}
		// Get last price for this store
		row := r.db.QueryRowContext(ctx,
			`SELECT price_cents FROM price_logs
			 WHERE user_id = ? AND LOWER(item_name) = LOWER(?) AND store = ?
			 ORDER BY date DESC LIMIT 1`,
			userID, itemName, s.Store,
		)
		row.Scan(&s.LastPrice)
		stores = append(stores, s)
	}
	return stores, nil
}

func (r *sqlitePriceRepo) GetSpendingSummary(ctx context.Context, userID int64, month string) (*domain.SpendingSummary, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(price_cents), 0), COUNT(*)
		 FROM price_logs
		 WHERE user_id = ? AND date LIKE ?`,
		userID, month+"%",
	)
	var s domain.SpendingSummary
	s.Month = month
	if err := row.Scan(&s.TotalCents, &s.ItemCount); err != nil {
		return nil, fmt.Errorf("spending summary: %w", err)
	}
	return &s, nil
}
