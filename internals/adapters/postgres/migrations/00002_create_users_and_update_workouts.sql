-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Existing workouts table migration
ALTER TABLE workouts
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id TYPE UUID USING gen_random_uuid(),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE workouts RENAME COLUMN name TO title;

ALTER TABLE workouts
    ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS date DATE,
    ADD COLUMN IF NOT EXISTS notes TEXT;

ALTER TABLE workouts DROP COLUMN IF EXISTS updated_at;

-- +goose Down
ALTER TABLE workouts DROP COLUMN IF EXISTS notes;
ALTER TABLE workouts DROP COLUMN IF EXISTS date;
ALTER TABLE workouts DROP COLUMN IF EXISTS user_id;
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE workouts RENAME COLUMN title TO name;

DROP TABLE IF EXISTS users;