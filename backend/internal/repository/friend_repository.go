package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type FriendRepository interface {
	CreateInvite(ctx context.Context, inv *domain.FriendInvite) error
	GetInviteByID(ctx context.Context, id int64) (*domain.FriendInvite, error)
	GetPendingInvites(ctx context.Context, email string) ([]domain.FriendInvite, error)
	AcceptInvite(ctx context.Context, inviteID, userID int64) error
	DeclineInvite(ctx context.Context, inviteID, userID int64) error
	AddFriendship(ctx context.Context, userID, friendID int64) error
	GetFriends(ctx context.Context, userID int64) ([]domain.FriendInfo, error)
	RemoveFriend(ctx context.Context, userID, friendID int64) error
	GetLeaderboard(ctx context.Context, userID int64) ([]domain.LeaderboardEntry, error)
}

type sqliteFriendRepo struct {
	db *sql.DB
}

func NewSQLiteFriendRepository(db *sql.DB) FriendRepository {
	return &sqliteFriendRepo{db: db}
}

func (r *sqliteFriendRepo) CreateInvite(ctx context.Context, inv *domain.FriendInvite) error {
	// Check if target email exists as a user
	var toUserID sql.NullInt64
	r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = ?`, inv.ToEmail).Scan(&toUserID)
	if toUserID.Valid {
		inv.ToUserID = &toUserID.Int64
	}

	res, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO friend_invites (from_user_id, to_email, to_user_id, status)
		 VALUES (?, ?, ?, 'pending')`,
		inv.FromUserID, inv.ToEmail, inv.ToUserID,
	)
	if err != nil {
		return fmt.Errorf("create invite: %w", err)
	}
	id, _ := res.LastInsertId()
	inv.ID = id
	return nil
}

func (r *sqliteFriendRepo) GetInviteByID(ctx context.Context, id int64) (*domain.FriendInvite, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, from_user_id, to_email, to_user_id, status, created_at
		 FROM friend_invites WHERE id = ?`, id,
	)
	var inv domain.FriendInvite
	var toUserID sql.NullInt64
	err := row.Scan(&inv.ID, &inv.FromUserID, &inv.ToEmail, &toUserID, &inv.Status, &inv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get invite: %w", err)
	}
	if toUserID.Valid {
		inv.ToUserID = &toUserID.Int64
	}
	return &inv, nil
}

func (r *sqliteFriendRepo) GetPendingInvites(ctx context.Context, email string) ([]domain.FriendInvite, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, from_user_id, to_email, to_user_id, status, created_at
		 FROM friend_invites WHERE to_email = ? AND status = 'pending'
		 ORDER BY created_at DESC`, email,
	)
	if err != nil {
		return nil, fmt.Errorf("query invites: %w", err)
	}
	defer rows.Close()

	var invites []domain.FriendInvite
	for rows.Next() {
		var inv domain.FriendInvite
		var toUserID sql.NullInt64
		if err := rows.Scan(&inv.ID, &inv.FromUserID, &inv.ToEmail, &toUserID, &inv.Status, &inv.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan invite: %w", err)
		}
		if toUserID.Valid {
			inv.ToUserID = &toUserID.Int64
		}
		invites = append(invites, inv)
	}
	return invites, nil
}

func (r *sqliteFriendRepo) AcceptInvite(ctx context.Context, inviteID, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE friend_invites SET status = 'accepted', to_user_id = ?
		 WHERE id = ? AND (to_user_id = ? OR to_user_id IS NULL)`,
		userID, inviteID, userID,
	)
	if err != nil {
		return fmt.Errorf("accept invite: %w", err)
	}
	return nil
}

func (r *sqliteFriendRepo) DeclineInvite(ctx context.Context, inviteID, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE friend_invites SET status = 'declined'
		 WHERE id = ? AND (to_user_id = ? OR to_user_id IS NULL)`,
		inviteID, userID,
	)
	if err != nil {
		return fmt.Errorf("decline invite: %w", err)
	}
	return nil
}

func (r *sqliteFriendRepo) AddFriendship(ctx context.Context, userID, friendID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Bidirectional friendship
	_, err = tx.ExecContext(ctx,
		`INSERT OR IGNORE INTO friendships (user_id, friend_id) VALUES (?, ?)`,
		userID, friendID,
	)
	if err != nil {
		return fmt.Errorf("add friendship a->b: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`INSERT OR IGNORE INTO friendships (user_id, friend_id) VALUES (?, ?)`,
		friendID, userID,
	)
	if err != nil {
		return fmt.Errorf("add friendship b->a: %w", err)
	}

	return tx.Commit()
}

func (r *sqliteFriendRepo) GetFriends(ctx context.Context, userID int64) ([]domain.FriendInfo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT u.id, u.name, u.email, f.created_at
		 FROM friendships f
		 JOIN users u ON u.id = f.friend_id
		 WHERE f.user_id = ?
		 ORDER BY u.name ASC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query friends: %w", err)
	}
	defer rows.Close()

	var friends []domain.FriendInfo
	for rows.Next() {
		var f domain.FriendInfo
		if err := rows.Scan(&f.UserID, &f.Name, &f.Email, &f.AddedAt); err != nil {
			return nil, fmt.Errorf("scan friend: %w", err)
		}
		friends = append(friends, f)
	}
	return friends, nil
}

func (r *sqliteFriendRepo) RemoveFriend(ctx context.Context, userID, friendID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM friendships WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)`,
		userID, friendID, friendID, userID,
	)
	if err != nil {
		return fmt.Errorf("remove friend: %w", err)
	}
	return nil
}

func (r *sqliteFriendRepo) GetLeaderboard(ctx context.Context, userID int64) ([]domain.LeaderboardEntry, error) {
	// Get current user + friends, rank by this week's calories
	rows, err := r.db.QueryContext(ctx,
		`WITH friend_ids AS (
			SELECT friend_id AS uid FROM friendships WHERE user_id = ?
			UNION SELECT ? AS uid
		),
		week_stats AS (
			SELECT fl.user_id, COALESCE(SUM(fl.calories), 0) AS week_cal
			FROM food_logs fl
			WHERE fl.user_id IN (SELECT uid FROM friend_ids)
			  AND fl.date >= date('now', 'weekday 0', '-6 days')
			GROUP BY fl.user_id
		)
		SELECT u.id, u.name, COALESCE(ws.week_cal, 0)
		FROM users u
		JOIN friend_ids fi ON fi.uid = u.id
		LEFT JOIN week_stats ws ON ws.user_id = u.id
		ORDER BY COALESCE(ws.week_cal, 0) DESC`,
		userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []domain.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e domain.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Name, &e.WeekCalories); err != nil {
			return nil, fmt.Errorf("scan leaderboard: %w", err)
		}
		e.Rank = rank
		rank++
		entries = append(entries, e)
	}
	return entries, nil
}
