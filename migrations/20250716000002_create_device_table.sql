-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    min_width INTEGER,
    max_width INTEGER,
    min_height INTEGER,
    max_height INTEGER,
    aspect_ratio_tolerance REAL,
    nsfw_mode INTEGER NOT NULL DEFAULT 0, -- 0=unspecified, 1=block, 2=accept, 3=only
    min_file_size INTEGER,
    max_file_size INTEGER,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS device_sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id INTEGER NOT NULL,
    source_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
    UNIQUE(device_id, source_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_devices_name ON devices(name);
CREATE INDEX IF NOT EXISTS idx_devices_slug ON devices(slug);
CREATE INDEX IF NOT EXISTS idx_devices_nsfw_mode ON devices(nsfw_mode);
CREATE INDEX IF NOT EXISTS idx_device_sources_device_id ON device_sources(device_id);
CREATE INDEX IF NOT EXISTS idx_device_sources_source_id ON device_sources(source_id);

-- +goose Down
DROP TABLE IF EXISTS device_sources;
DROP TABLE IF EXISTS devices;