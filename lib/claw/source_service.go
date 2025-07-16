package claw

import "database/sql"

// SourceService provides business logic for managing sources
type SourceService struct {
	db *sql.DB
}

// NewSourceService creates a new SourceService
func NewSourceService(db *sql.DB) *SourceService {
	return &SourceService{db: db}
}