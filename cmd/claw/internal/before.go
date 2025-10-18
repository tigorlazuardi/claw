package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/tigorlazuardi/claw/lib/logger"
	"github.com/tigorlazuardi/claw/lib/otel"
	"github.com/urfave/cli/v3"
)

var cwd, _ = os.Getwd()

// Before is a CLI hook that runs before any command. It initializes and watches the configuration file.
func Before(ctx context.Context, c *cli.Command) (context.Context, error) {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	if loc, err := time.LoadLocation(tz); err == nil {
		time.Local = loc
	}
	if err := logger.Setup(ctx); err != nil {
		return ctx, fmt.Errorf("failed to setup logger: %w", err)
	}
	if err := otel.SetupTracing(ctx); err != nil {
		return ctx, fmt.Errorf("failed to setup OpenTelemetry tracing: %w", err)
	}

	if cfg == nil {
		cfg = defaultCLIConfig()
	}
	if err := cfg.ReadAndWatch(); err != nil {
		return ctx, fmt.Errorf("failed to read config and watch file: %w", err)
	}
	if err := cfg.LoadEnv(); err != nil {
		return ctx, fmt.Errorf("failed to read environment variable into config: %w", err)
	}
	slog.Info("watching config file", "path", cfg.FileLocation)
	return ctx, nil
}

// ConfigFlag returns a CLI flag for specifying the configuration file path.
func ConfigFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "config",
		Aliases:  []string{"c"},
		Required: false,
		Usage:    "Path to configuration file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CLAW_CONFIG"),
		),
		Value: func() string {
			return filepath.Join(xdg.ConfigHome, "claw", "config.yaml")
		}(),
		Validator: func(s string) error {
			if s == "" {
				return errors.New("config path cannot be empty")
			}
			info, err := os.Stat(s)
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to stat config file: %w", err)
			}
			if info.IsDir() {
				return fmt.Errorf("config path %q is a directory", s)
			}
			return nil
		},
	}
}
