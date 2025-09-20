-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    is_disabled INTEGER NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    slug TEXT NOT NULL UNIQUE,
    save_dir TEXT NOT NULL DEFAULT '',
    filename_template TEXT NOT NULL DEFAULT '',
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    aspect_ratio_difference REAL NOT NULL DEFAULT 0.0,
    image_min_width INTEGER NOT NULL DEFAULT 0,
    image_max_width INTEGER NOT NULL DEFAULT 0,
    image_min_height INTEGER NOT NULL DEFAULT 0,
    image_max_height INTEGER NOT NULL DEFAULT 0,
    image_min_file_size INTEGER NOT NULL DEFAULT 0,
    image_max_file_size INTEGER NOT NULL DEFAULT 0,
    nsfw_mode INTEGER NOT NULL DEFAULT 0, -- 0=unspecified, 1=allow, 2=block, 3=only
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    last_active_at INTEGER -- Last time an image was downloaded for this device
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
CREATE INDEX IF NOT EXISTS idx_devices_is_disabled ON devices(is_disabled);
CREATE INDEX IF NOT EXISTS idx_device_sources_device_id ON device_sources(device_id);
CREATE INDEX IF NOT EXISTS idx_device_sources_source_id ON device_sources(source_id);

-- +goose Down
DROP TABLE IF EXISTS device_sources;
DROP TABLE IF EXISTS devices;
