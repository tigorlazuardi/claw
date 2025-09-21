-- +goose Up
CREATE TABLE IF NOT EXISTS images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id INTEGER NOT NULL,
    download_url TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    filesize INTEGER NOT NULL,
    thumbnail_path TEXT NOT NULL,
    image_path TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    post_author TEXT NOT NULL DEFAULT '',
    post_author_url TEXT NOT NULL DEFAULT '',
    post_url TEXT NOT NULL DEFAULT '',
    is_favorite INTEGER NOT NULL DEFAULT 0, -- 0=false, 1=true
    is_nsfw INTEGER NOT NULL DEFAULT 0, -- 0=false, 1=true
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_images_source_id ON images(source_id);
CREATE INDEX IF NOT EXISTS idx_images_is_favorite ON images(is_favorite);
CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at);
CREATE INDEX IF NOT EXISTS idx_images_updated_at ON images(updated_at);

-- Junction table for image-device many-to-many relationship
-- and contain path to the image on that device.
CREATE TABLE IF NOT EXISTS image_devices (
    image_id INTEGER NOT NULL,
    device_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    path TEXT NOT NULL,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    PRIMARY KEY (image_id, device_id),
    UNIQUE(path)
);

CREATE INDEX IF NOT EXISTS idx_image_devices_device_id_image_id ON image_devices(device_id, image_id); -- reverse composite index if needed to search by device_id first.

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_image_devices_on_insert_update_device_last_active
AFTER INSERT ON image_devices
FOR EACH ROW
BEGIN
  UPDATE devices 
  SET last_active_at = unixepoch('now', 'subsec') * 1000
  WHERE id = NEW.device_id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS image_devices;
DROP TABLE IF EXISTS images;
