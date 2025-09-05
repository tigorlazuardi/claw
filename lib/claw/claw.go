package claw

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"sync"

	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/claw/lib/claw/config"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	"github.com/tigorlazuardi/claw/lib/claw/source"
	"github.com/tigorlazuardi/claw/lib/claw/source/reddit"
	"golang.org/x/sync/semaphore"
)

// Option is a function that configures a Claw instance
type Option func(*Claw)

// WithLogger sets a custom logger for the Claw instance
func WithLogger(logger *slog.Logger) Option {
	return func(c *Claw) {
		c.logger = logger
	}
}

// WithBackend registers a source backend with the given name
func WithBackend(name string, backend source.Source) Option {
	return func(c *Claw) {
		if c.scheduler.backends == nil {
			c.scheduler.backends = make(map[string]source.Source)
		}
		c.scheduler.backends[name] = backend
	}
}

// WithHTTPClient sets a custom HTTP client for the scheduler
func WithHTTPClient(client *http.Client) Option {
	return func(c *Claw) {
		c.scheduler.httpclient = client
	}
}

// Claw provides business logic for managing sources
type Claw struct {
	db     *sql.DB
	config *config.Config
	logger *slog.Logger

	scheduler *scheduler
}

// New creates a new Claw instance
func New(db *sql.DB, config *config.Config, opts ...Option) *Claw {
	cl := &Claw{
		db:     db,
		logger: slog.Default(),
		config: config,
	}

	// Initialize scheduler with default backends
	cl.scheduler = &scheduler{
		claw:           cl,
		config:         config,
		queue:          make(chan model.Jobs, 1024),
		tracker:        &tracker{},
		imageSemaphore: semaphore.NewWeighted(leastCommonMultiple),
		wg:             &sync.WaitGroup{},
		reloadSignal:   broadcast.NewRelay[struct{}](),
		logger:         cl.logger,
		backends: map[string]source.Source{
			reddit.SourceName: &reddit.Reddit{
				Client: http.DefaultClient,
			},
		},
		httpclient: http.DefaultClient,
	}

	// Apply options
	for _, opt := range opts {
		opt(cl)
	}

	return cl
}

func (claw *Claw) RereadConfig() {
	claw.scheduler.reloadSignal.Broadcast(struct{}{})
}

// StartSchedculer starts the job scheduler
//
// It blocks until the given context is cancelled
func (claw *Claw) StartSchedculer(ctx context.Context) {
	claw.scheduler.start(ctx)
}
