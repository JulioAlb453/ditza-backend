ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS description VARCHAR(160) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS category VARCHAR(40) NOT NULL DEFAULT 'general',
    ADD COLUMN IF NOT EXISTS color VARCHAR(24) NOT NULL DEFAULT 'green',
    ADD COLUMN IF NOT EXISTS frequency VARCHAR(20) NOT NULL DEFAULT 'daily',
    ADD COLUMN IF NOT EXISTS target_count INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS target_unit VARCHAR(30) NOT NULL DEFAULT 'veces',
    ADD COLUMN IF NOT EXISTS difficulty VARCHAR(20) NOT NULL DEFAULT 'medium',
    ADD COLUMN IF NOT EXISTS reminder_time VARCHAR(5) NULL;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_habits_frequency') THEN
        ALTER TABLE habits
            ADD CONSTRAINT chk_habits_frequency
            CHECK (frequency IN ('daily', 'weekly', 'specific_days')) NOT VALID;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_habits_target_count') THEN
        ALTER TABLE habits
            ADD CONSTRAINT chk_habits_target_count
            CHECK (target_count > 0) NOT VALID;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_habits_difficulty') THEN
        ALTER TABLE habits
            ADD CONSTRAINT chk_habits_difficulty
            CHECK (difficulty IN ('easy', 'medium', 'hard')) NOT VALID;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_habits_reminder_time') THEN
        ALTER TABLE habits
            ADD CONSTRAINT chk_habits_reminder_time
            CHECK (reminder_time IS NULL OR reminder_time ~ '^([01][0-9]|2[0-3]):[0-5][0-9]$') NOT VALID;
    END IF;
END $$;

ALTER TABLE habits
    VALIDATE CONSTRAINT chk_habits_frequency,
    VALIDATE CONSTRAINT chk_habits_target_count,
    VALIDATE CONSTRAINT chk_habits_difficulty,
    VALIDATE CONSTRAINT chk_habits_reminder_time;
