package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type PantryRepository interface {
	SetItems(ctx context.Context, userID int64, items []string) error
	GetItems(ctx context.Context, userID int64) ([]domain.PantryItem, error)
}

type SQLitePantryRepository struct {
	db *sql.DB
}

func NewSQLitePantryRepository(db *sql.DB) *SQLitePantryRepository {
	return &SQLitePantryRepository{db: db}
}

func (r *SQLitePantryRepository) SetItems(ctx context.Context, userID int64, items []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM pantry_items WHERE user_id = ?`, userID); err != nil {
		return fmt.Errorf("clearing pantry: %w", err)
	}

	now := time.Now().UTC()
	for _, name := range items {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO pantry_items (user_id, name, created_at) VALUES (?, ?, ?)`,
			userID, name, now,
		); err != nil {
			return fmt.Errorf("inserting pantry item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *SQLitePantryRepository) GetItems(ctx context.Context, userID int64) ([]domain.PantryItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, name, created_at FROM pantry_items WHERE user_id = ? ORDER BY name`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying pantry items: %w", err)
	}
	defer rows.Close()

	var items []domain.PantryItem
	for rows.Next() {
		var item domain.PantryItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.Name, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning pantry item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
