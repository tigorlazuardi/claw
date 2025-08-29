package claw

import (
	"database/sql"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	"github.com/tigorlazuardi/claw/lib/claw/source"
)

// SchedulerConfig holds scheduler configuration
type SchedulerConfig struct {
	// PollInterval is how often to poll for new jobs (default: 5 seconds)
	PollInterval time.Duration
	// MaxWorkers is the maximum number of concurrent job workers (default: 3)
	MaxWorkers int
	// DownloadWorkers is the number of concurrent download workers (default: 5)
	DownloadWorkers int
	// BaseDir is the base directory for storing images
	BaseDir string
	// TmpDir is the temporary directory for downloads
	TmpDir string
}

// DefaultSchedulerConfig returns a default scheduler configuration
func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		PollInterval:    5 * time.Second,
		MaxWorkers:      3,
		DownloadWorkers: 5,
		BaseDir:         "./data",
		TmpDir:          os.TempDir(),
	}
}

// Claw provides business logic for managing sources
type Claw struct {
	db     *sql.DB
	logger *slog.Logger

	// Scheduler fields
	schedulerConfig  SchedulerConfig
	sources          map[string]source.Source
	jobQueue         chan *model.Jobs
	queuedJobs       map[int64]bool
	queuedJobsMutex  sync.RWMutex
	downloadQueue    chan downloadTask
	schedulerStopCh  chan struct{}
	schedulerDoneCh  chan struct{}
	schedulerRunning atomic.Bool
	schedulerMutex   sync.RWMutex
}

// downloadTask represents a single image download task
type downloadTask struct {
	jobID      int64
	sourceID   int64
	image      source.Image
	devices    []deviceFilter
	sourceName string
}

// deviceFilter contains device information for filtering
type deviceFilter struct {
	id                    int64
	slug                  string
	saveDir               string
	width                 int64
	height                int64
	aspectRatioDifference float64
	imageMinWidth         int64
	imageMaxWidth         int64
	imageMinHeight        int64
	imageMaxHeight        int64
	imageMinFileSize      int64
	imageMaxFileSize      int64
	nsfwMode              int64
}

// New creates a new Claw instance
func New(db *sql.DB) *Claw {
	return &Claw{
		db:              db,
		schedulerConfig: DefaultSchedulerConfig(),
		sources:         make(map[string]source.Source),
		queuedJobs:      make(map[int64]bool),
	}
}

// SetSchedulerConfig updates the scheduler configuration
func (c *Claw) SetSchedulerConfig(config SchedulerConfig) {
	c.schedulerConfig = config
}

// RegisterSource registers a source backend for the scheduler
func (c *Claw) RegisterSource(src source.Source) {
	if c.sources == nil {
		c.sources = make(map[string]source.Source)
	}
	c.sources[src.Name()] = src
}
