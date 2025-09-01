package claw

import (
	"database/sql"
	"log/slog"
	"sync"

	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/claw/lib/claw/config"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	"golang.org/x/sync/semaphore"
)

// Claw provides business logic for managing sources
type Claw struct {
	db     *sql.DB
	config *config.Config
	logger *slog.Logger

	scheduler *scheduler
}

// New creates a new Claw instance
func New(db *sql.DB, config *config.Config) *Claw {
	cl := &Claw{
		db:     db,
		logger: slog.Default().With("package", "github.com/tigorlazuardi/claw/lib/claw"),
		config: config,
	}
	cl.scheduler = &scheduler{
		claw:           cl,
		config:         config,
		queue:          make(chan model.Jobs, 1024),
		tracker:        &tracker{},
		imageSemaphore: semaphore.NewWeighted(leastCommonMultiple),
		wg:             &sync.WaitGroup{},
		reloadSignal:   broadcast.NewRelay[struct{}](),
		logger:         cl.logger.With("component", "scheduler"),
	}
	return cl
}

func (claw *Claw) RereadConfig() {
	claw.scheduler.reloadSignal.Broadcast(struct{}{})
}
