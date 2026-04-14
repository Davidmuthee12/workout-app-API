-- +goose Up
CREATE TABLE IF NOT EXISTS programs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	program_id UUID NOT NULL,
	title TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS services (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	service_id UUID NOT NULL,
	title TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS programs;
