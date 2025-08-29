package claw

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/tigorlazuardi/claw/lib/claw/source/reddit"
)

// ExampleSchedulerUsage demonstrates how to use the integrated scheduler
func ExampleSchedulerUsage(claw *Claw) {
	// Configure logger if not already set
	if claw.logger == nil {
		claw.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	// Register source backends
	claw.RegisterSource(&reddit.Reddit{Client: &http.Client{}})

	// Configure scheduler
	config := DefaultSchedulerConfig()
	config.BaseDir = "./data"
	config.TmpDir = "/tmp"
	claw.SetSchedulerConfig(config)

	// Start scheduler in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go claw.StartScheduler(ctx)

	// In real usage, the scheduler would run indefinitely
	// To stop the scheduler gracefully:
	// claw.StopScheduler()
	// claw.WaitScheduler()
}

