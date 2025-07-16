-- +goose Up
CREATE TABLE IF NOT EXISTS sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL,
    display_name TEXT NOT NULL,
    parameter TEXT NOT NULL,
    countback INTEGER NOT NULL DEFAULT 0,
    last_run_at INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    UNIQUE(slug, parameter)
);

CREATE TABLE IF NOT EXISTS schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id INTEGER NOT NULL,
    slug_id INTEGER NOT NULL,
    schedule TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
    FOREIGN KEY (slug_id) REFERENCES sources(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_sources_slug ON sources(slug);
CREATE INDEX IF NOT EXISTS idx_sources_last_run_at ON sources(last_run_at);
CREATE INDEX IF NOT EXISTS idx_sources_slug_parameter ON sources(slug, parameter);
CREATE INDEX IF NOT EXISTS idx_schedules_source_id ON schedules(source_id);
CREATE INDEX IF NOT EXISTS idx_schedules_slug_id ON schedules(slug_id);

-- +goose Down
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS sources;
