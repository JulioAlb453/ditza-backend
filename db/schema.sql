CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alias VARCHAR(60) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    coins INTEGER NOT NULL DEFAULT 0 CHECK (coins >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS habits (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(80) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    current_streak INTEGER NOT NULL DEFAULT 0 CHECK (current_streak >= 0),
    best_streak INTEGER NOT NULL DEFAULT 0 CHECK (best_streak >= 0),
    last_completed_date DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS habit_completions (
    id BIGSERIAL PRIMARY KEY,
    habit_id BIGINT NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    completed_at TIMESTAMPTZ NOT NULL,
    -- AT TIME ZONE 'UTC' hace la expresión IMMUTABLE (requerido en columnas GENERATED)
    completion_date DATE GENERATED ALWAYS AS ((completed_at AT TIME ZONE 'UTC')::date) STORED,
    note VARCHAR(140) NULL,
    emoji VARCHAR(16) NULL,
    coins_awarded INTEGER NOT NULL,
    season_points_awarded INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_habit_completion_per_day UNIQUE (habit_id, completion_date)
);

CREATE TABLE IF NOT EXISTS cosmetics (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    slot VARCHAR(20) NOT NULL CHECK (slot IN ('hat', 'shirt', 'background', 'accessory')),
    price_coins INTEGER NOT NULL CHECK (price_coins > 0),
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('common', 'rare')),
    asset_key VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pets (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(30) NOT NULL,
    level INTEGER NOT NULL DEFAULT 1 CHECK (level >= 1),
    xp INTEGER NOT NULL DEFAULT 0 CHECK (xp >= 0),
    mood VARCHAR(16) NOT NULL CHECK (mood IN ('happy', 'neutral', 'sad', 'sleeping')),
    equipped_hat_id BIGINT NULL REFERENCES cosmetics(id) ON DELETE SET NULL,
    equipped_shirt_id BIGINT NULL REFERENCES cosmetics(id) ON DELETE SET NULL,
    equipped_background_id BIGINT NULL REFERENCES cosmetics(id) ON DELETE SET NULL,
    equipped_accessory_id BIGINT NULL REFERENCES cosmetics(id) ON DELETE SET NULL,
    last_interaction_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_cosmetics (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cosmetic_id BIGINT NOT NULL REFERENCES cosmetics(id) ON DELETE CASCADE,
    purchased_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, cosmetic_id)
);

CREATE TABLE IF NOT EXISTS friendships (
    id BIGSERIAL PRIMARY KEY,
    requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    addressee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    responded_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_friendship_different_users CHECK (requester_id <> addressee_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_friendships_pair
    ON friendships (LEAST(requester_id, addressee_id), GREATEST(requester_id, addressee_id));

CREATE TABLE IF NOT EXISTS seasons (
    id BIGSERIAL PRIMARY KEY,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_season_dates CHECK (ends_at > starts_at)
);

CREATE TABLE IF NOT EXISTS season_scores (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    season_id BIGINT NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    points INTEGER NOT NULL DEFAULT 0 CHECK (points >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, season_id)
);

CREATE TABLE IF NOT EXISTS point_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL CHECK (type IN ('habit', 'streak_bonus', 'note_bonus', 'purchase', 'season_reset')),
    coins_delta INTEGER NOT NULL DEFAULT 0,
    season_delta INTEGER NOT NULL DEFAULT 0,
    reference_id BIGINT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_point_delta_non_zero CHECK (coins_delta <> 0 OR season_delta <> 0)
);

CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits(user_id);
CREATE INDEX IF NOT EXISTS idx_habits_active ON habits(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_habit_completions_user_date ON habit_completions(user_id, completion_date);
CREATE INDEX IF NOT EXISTS idx_habit_completions_habit_id ON habit_completions(habit_id);
CREATE INDEX IF NOT EXISTS idx_friendships_requester ON friendships(requester_id);
CREATE INDEX IF NOT EXISTS idx_friendships_addressee ON friendships(addressee_id);
CREATE INDEX IF NOT EXISTS idx_friendships_status ON friendships(status);
CREATE INDEX IF NOT EXISTS idx_seasons_active ON seasons(is_active);
CREATE INDEX IF NOT EXISTS idx_season_scores_season_points ON season_scores(season_id, points DESC);
CREATE INDEX IF NOT EXISTS idx_point_transactions_user_created ON point_transactions(user_id, created_at DESC);

-- Semilla mínima de temporada activa (15 días)
INSERT INTO seasons (starts_at, ends_at, is_active)
SELECT NOW(), NOW() + INTERVAL '15 days', TRUE
WHERE NOT EXISTS (SELECT 1 FROM seasons WHERE is_active = TRUE);

-- Semilla mínima de cosméticos
INSERT INTO cosmetics (name, slot, price_coins, rarity, asset_key, is_active)
SELECT *
FROM (
    VALUES
        ('Gorra Verde', 'hat', 120, 'common', 'hat_green_cap', TRUE),
        ('Camiseta Azul', 'shirt', 150, 'common', 'shirt_blue_basic', TRUE),
        ('Fondo Bosque', 'background', 220, 'rare', 'bg_forest_01', TRUE),
        ('Lentes Retro', 'accessory', 180, 'common', 'acc_retro_glasses', TRUE)
) AS seed(name, slot, price_coins, rarity, asset_key, is_active)
WHERE NOT EXISTS (SELECT 1 FROM cosmetics);
