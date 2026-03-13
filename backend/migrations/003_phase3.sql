-- BiteWise Phase 3: Achievements, Custom Recipes, Sharing, Notifications, Stats

-- Seed default achievements
INSERT OR IGNORE INTO achievements (name, description, icon, category, condition_type, condition_value) VALUES
    ('first_log', 'Erstes Lebensmittel erfasst', 'nutrition', 'milestone', 'food_logs_total', 1),
    ('streak_3', '3 Tage in Folge getrackt', 'flame', 'streak', 'tracking_streak', 3),
    ('streak_7', '7 Tage in Folge getrackt', 'flame', 'streak', 'tracking_streak', 7),
    ('streak_30', '30 Tage in Folge getrackt', 'flame', 'streak', 'tracking_streak', 30),
    ('water_champion', 'Wasserziel 7 Tage in Folge erreicht', 'water', 'streak', 'water_streak', 7),
    ('recipe_explorer', '10 verschiedene Rezepte ausprobiert', 'restaurant', 'milestone', 'unique_recipes', 10),
    ('meal_planner', '5 Essenspläne generiert', 'calendar', 'milestone', 'meal_plans_generated', 5),
    ('barcode_scanner', 'Ersten Barcode gescannt', 'barcode', 'milestone', 'barcodes_scanned', 1),
    ('pantry_chef', '3 Rezepte aus Vorrats-Vorschlägen gekocht', 'leaf', 'milestone', 'pantry_recipes', 3),
    ('seasonal_eater', '5 saisonale Zutaten in einem Monat erfasst', 'sunny', 'milestone', 'seasonal_ingredients', 5),
    ('healthy_week', 'Kalorienziel ±10% für 7 Tage eingehalten', 'heart', 'streak', 'calorie_target_streak', 7),
    ('macro_master', 'Alle Makro-Ziele an einem Tag erreicht', 'trophy', 'milestone', 'macro_targets_hit', 1);

-- Add user_id to recipes (NULL = system/scraped, set = user-created)
ALTER TABLE recipes ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;

-- Shared recipes
CREATE TABLE IF NOT EXISTS shared_recipes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    share_code TEXT NOT NULL UNIQUE,
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shared_recipes_code ON shared_recipes(share_code);

-- User notifications / reminders
CREATE TABLE IF NOT EXISTS user_notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('meal_reminder_breakfast', 'meal_reminder_lunch', 'meal_reminder_dinner', 'water_reminder', 'weekly_plan_reminder')),
    time TEXT NOT NULL DEFAULT '08:00',
    enabled INTEGER NOT NULL DEFAULT 1,
    UNIQUE(user_id, type)
);

CREATE INDEX IF NOT EXISTS idx_user_notifications_user ON user_notifications(user_id);
