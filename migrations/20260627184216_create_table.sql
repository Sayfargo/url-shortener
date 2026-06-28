-- +goose Up
CREATE TABLE IF NOT EXISTS shorted_urls (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    original_url TEXT,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS shorted_urls;
