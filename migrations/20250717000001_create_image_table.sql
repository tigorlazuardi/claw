-- +goose Up
CREATE TABLE IF NOT EXISTS images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id INTEGER NOT NULL,
    download_url TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    filesize INTEGER NOT NULL,
    thumbnail_path TEXT,
    image_path TEXT NOT NULL,
    post_author TEXT,
    post_author_url TEXT,
    post_url TEXT,
    is_favorite INTEGER NOT NULL DEFAULT 0, -- 0=false, 1=true
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE
);

-- Junction table for image-device many-to-many relationship
CREATE TABLE IF NOT EXISTS image_devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image_id INTEGER NOT NULL,
    device_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    UNIQUE(image_id, device_id)
);

-- Table for image file paths (hardlinks/copies)
CREATE TABLE IF NOT EXISTS image_paths (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

-- Table for image tags
CREATE TABLE IF NOT EXISTS image_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    UNIQUE(image_id, tag)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_images_source_id ON images(source_id);
CREATE INDEX IF NOT EXISTS idx_images_download_url ON images(download_url);
CREATE INDEX IF NOT EXISTS idx_images_is_favorite ON images(is_favorite);
CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at);
CREATE INDEX IF NOT EXISTS idx_images_updated_at ON images(updated_at);

CREATE INDEX IF NOT EXISTS idx_image_devices_image_id ON image_devices(image_id);
CREATE INDEX IF NOT EXISTS idx_image_devices_device_id ON image_devices(device_id);

CREATE INDEX IF NOT EXISTS idx_image_paths_image_id ON image_paths(image_id);
CREATE INDEX IF NOT EXISTS idx_image_paths_path ON image_paths(path);

CREATE INDEX IF NOT EXISTS idx_image_tags_image_id ON image_tags(image_id);
CREATE INDEX IF NOT EXISTS idx_image_tags_tag ON image_tags(tag);

-- +goose Down
DROP TABLE IF EXISTS image_tags;
DROP TABLE IF EXISTS image_paths;
DROP TABLE IF EXISTS image_devices;
DROP TABLE IF EXISTS images;