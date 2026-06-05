UPDATE habits
SET reminder_time = ''
WHERE reminder_time IS NULL;

ALTER TABLE habits
    ALTER COLUMN reminder_time SET DEFAULT '',
    ALTER COLUMN reminder_time SET NOT NULL;

ALTER TABLE habits
    DROP CONSTRAINT IF EXISTS chk_habits_reminder_time;

ALTER TABLE habits
    ADD CONSTRAINT chk_habits_reminder_time
    CHECK (reminder_time = '' OR reminder_time ~ '^([01][0-9]|2[0-3]):[0-5][0-9]$');
