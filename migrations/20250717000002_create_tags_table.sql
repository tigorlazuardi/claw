-- +goose Up
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at INTEGER NOT NULL
);

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);

-- +goose Down
DROP TABLE IF EXISTS tags;