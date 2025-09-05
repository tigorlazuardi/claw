package config

import (
	"log/slog"
	"time"
)

type Scheduler struct {
	// PollInterval is how often to poll for new jobs (default: 5 seconds).
	PollInterval time.Duration `koanf:"poll_interval"`
	// MaxWorkers is the maximum number of concurrent jobs that can be processed at the same time (default: 3).
	//
	// Value is clamped between 1 and 16.
	MaxWorkers int `koanf:"max_workers"`
	// DownloadWorkers is the number of concurrent download workers (default: 5)
	//
	// Value is clamped between 1 and 16.
	DownloadWorkers int `koanf:"download_workers"`

	// ExitTimeout is the time to wait for workers to finish when shutting down (default: 10 seconds).
	ExitTimeout time.Duration `koanf:"exit_timeout"`
}

func (sc Scheduler) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("max_workers", sc.MaxWorkers),
		slog.Int("download_workers", sc.DownloadWorkers),
		slog.Duration("poll_interval", sc.PollInterval),
		slog.Duration("exit_timeout", sc.ExitTimeout),
	)
}

func DefaultScheduler() Scheduler {
	return Scheduler{
		PollInterval:    5 * time.Second,
		MaxWorkers:      3,
		DownloadWorkers: 5,
		ExitTimeout:     10 * time.Second,
	}
}
