-- +goose Up
-- +goose NO TRANSACTION
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA temp_store = MEMORY;
PRAGMA mmap_size = 268435456;
PRAGMA cache_size = 10000;

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd