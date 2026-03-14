package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type DiaryRepository interface {
	Create(ctx context.Context, entry *domain.DiaryEntry) error
	GetByDate(ctx context.Context, userID int64, date string) (*domain.DiaryEntry, error)
	GetByDateRange(ctx context.Context, userID int64, from, to string) ([]domain.DiaryEntry, error)
	Update(ctx context.Context, entry *domain.DiaryEntry) error
	Delete(ctx context.Context, userID, id int64) error
}

type sqliteDiaryRepo struct {
	db *sql.DB
}

func NewSQLiteDiaryRepository(db *sql.DB) DiaryRepository {
	return &sqliteDiaryRepo{db: db}
}

func (r *sqliteDiaryRepo) Create(ctx context.Context, entry *domain.DiaryEntry) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO diary_entries (user_id, date, mood, energy_level, notes)
		 VALUES (?, ?, ?, ?, ?)`,
		entry.UserID, entry.Date, entry.Mood, entry.EnergyLevel, entry.Notes,
	)
	if err != nil {
		return fmt.Errorf("insert diary entry: %w", err)
	}
	id, _ := res.LastInsertId()
	entry.ID = id
	return nil
}

func (r *sqliteDiaryRepo) GetByDate(ctx context.Context, userID int64, date string) (*domain.DiaryEntry, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, date, mood, energy_level, notes, created_at, updated_at
		 FROM diary_entries WHERE user_id = ? AND date = ?`, userID, date,
	)
	var e domain.DiaryEntry
	err := row.Scan(&e.ID, &e.UserID, &e.Date, &e.Mood, &e.EnergyLevel, &e.Notes, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get diary entry: %w", err)
	}
	return &e, nil
}

func (r *sqliteDiaryRepo) GetByDateRange(ctx context.Context, userID int64, from, to string) ([]domain.DiaryEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, date, mood, energy_level, notes, created_at, updated_at
		 FROM diary_entries WHERE user_id = ? AND date >= ? AND date <= ?
		 ORDER BY date ASC`, userID, from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("query diary entries: %w", err)
	}
	defer rows.Close()

	var entries []domain.DiaryEntry
	for rows.Next() {
		var e domain.DiaryEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Date, &e.Mood, &e.EnergyLevel, &e.Notes, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan diary entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *sqliteDiaryRepo) Update(ctx context.Context, entry *domain.DiaryEntry) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE diary_entries SET mood = ?, energy_level = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ? AND user_id = ?`,
		entry.Mood, entry.EnergyLevel, entry.Notes, entry.ID, entry.UserID,
	)
	if err != nil {
		return fmt.Errorf("update diary entry: %w", err)
	}
	return nil
}

func (r *sqliteDiaryRepo) Delete(ctx context.Context, userID, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM diary_entries WHERE id = ? AND user_id = ?`, id, userID,
	)
	if err != nil {
		return fmt.Errorf("delete diary entry: %w", err)
	}
	return nil
}
