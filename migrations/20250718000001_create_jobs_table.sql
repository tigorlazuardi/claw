-- +goose Up
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id INTEGER NOT NULL,
    schedule_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    run_at INTEGER,
    finished_at INTEGER,
    status TEXT NOT NULL DEFAULT 'pending',
    error TEXT,
    FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS job_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER NOT NULL,
    image_id INTEGER NOT NULL,
    device_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    UNIQUE(job_id, image_id, device_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_jobs_source_id ON jobs(source_id);
CREATE INDEX IF NOT EXISTS idx_jobs_schedule_id ON jobs(schedule_id);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_run_at ON jobs(run_at);
CREATE INDEX IF NOT EXISTS idx_jobs_finished_at ON jobs(finished_at);
CREATE INDEX IF NOT EXISTS idx_job_images_job_id ON job_images(job_id);
CREATE INDEX IF NOT EXISTS idx_job_images_image_id ON job_images(image_id);
CREATE INDEX IF NOT EXISTS idx_job_images_device_id ON job_images(device_id);
CREATE INDEX IF NOT EXISTS idx_job_images_action ON job_images(action);

-- +goose Down
DROP TABLE IF EXISTS job_images;
DROP TABLE IF EXISTS jobs;
