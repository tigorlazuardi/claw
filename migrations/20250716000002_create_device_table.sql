-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL COLLATE NOCASE,
    is_disabled INTEGER NOT NULL DEFAULT 0,
    name TEXT NOT NULL UNIQUE COLLATE NOCASE,
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
    nsfw_mode INTEGER NOT NULL DEFAULT 1, -- 0=unspecified, 1=allow, 2=block, 3=only
    created_at INTEGER NOT NULL DEFAULT 0,
    updated_at INTEGER NOT NULL DEFAULT 0,
    last_active_at INTEGER NOT NULL DEFAULT 0, -- Last time an image was downloaded for this device
    UNIQUE(slug)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_devices_name ON devices(name COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS idx_devices_nsfw_mode ON devices(nsfw_mode);
CREATE INDEX IF NOT EXISTS idx_devices_is_disabled ON devices(is_disabled);

-- landing page index is index when client requests images for landing page but
-- sets no filter and sorts. So the default page is optimized to display active devices first, then
-- the latest active devices.
CREATE INDEX IF NOT EXISTS idx_devices_landing_page ON devices(is_disabled, last_active_at);

CREATE TABLE IF NOT EXISTS device_sources (
    device_id INTEGER NOT NULL,
    source_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
    PRIMARY KEY (device_id, source_id)
);

CREATE INDEX IF NOT EXISTS idx_device_sources_source_id_device_id ON device_sources(source_id, device_id); -- reverse composite index if needed to search by source_id first.

-- +goose Down
DROP TABLE IF EXISTS device_sources;
DROP TABLE IF EXISTS devices;
