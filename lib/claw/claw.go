package claw

import "database/sql"

// Claw provides business logic for managing sources
type Claw struct {
	db *sql.DB
}

// New creates a new SourceService
func New(db *sql.DB) *Claw {
	return &Claw{db: db}
}

