package claw

import (
	"database/sql"
	"log/slog"
)

// Claw provides business logic for managing sources
type Claw struct {
	db     *sql.DB
	Logger *slog.Logger
}

// New creates a new SourceService
func New(db *sql.DB) *Claw {
	return &Claw{db: db}
}
