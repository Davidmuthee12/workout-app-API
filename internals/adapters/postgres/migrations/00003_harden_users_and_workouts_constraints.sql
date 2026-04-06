-- +goose Up
-- Ensure every workout is owned by a user before enforcing NOT NULL.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM workouts WHERE user_id IS NULL) THEN
        RAISE EXCEPTION 'cannot set workouts.user_id NOT NULL: existing rows have NULL user_id';
    END IF;
END $$;

ALTER TABLE workouts
    ALTER COLUMN user_id SET NOT NULL;

-- Keep timestamp semantics consistent across tables.
ALTER TABLE users
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

-- +goose Down
ALTER TABLE users
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'UTC';

ALTER TABLE workouts
    ALTER COLUMN user_id DROP NOT NULL;
