-- Phase 4 Migration: Diary, Prices, Friends

-- Food Diary / Journal
CREATE TABLE IF NOT EXISTS diary_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date TEXT NOT NULL,
    mood TEXT NOT NULL DEFAULT 'neutral', -- great, good, neutral, bad, terrible
    energy_level INTEGER NOT NULL DEFAULT 5, -- 1-10
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_diary_entries_user_date ON diary_entries(user_id, date);

-- Grocery Price Tracking
CREATE TABLE IF NOT EXISTS price_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_name TEXT NOT NULL,
    price_cents INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'EUR',
    store TEXT NOT NULL DEFAULT '',
    date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_price_logs_user ON price_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_price_logs_item ON price_logs(user_id, item_name);

-- Friendships
CREATE TABLE IF NOT EXISTS friendships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, friend_id)
);

CREATE TABLE IF NOT EXISTS friend_invites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_email TEXT NOT NULL,
    to_user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, accepted, declined
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_user_id, to_email)
);
CREATE INDEX IF NOT EXISTS idx_friend_invites_to_email ON friend_invites(to_email);
CREATE INDEX IF NOT EXISTS idx_friend_invites_to_user ON friend_invites(to_user_id);
