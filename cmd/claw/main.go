package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tigorlazuardi/claw/cmd/claw/internal"
	"github.com/urfave/cli/v3"
)

// Change the version file inside CI/CD pipeline during build time using go build -ldflags "-X main.Version=your_version"
var Version = "v0.0.0"

func main() {
	godotenv.Load()

	// Create context that cancels on interrupt/termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-ctx.Done()
		// Ensure we only listen for signals once. We want consecutive signals to force exit the
		// program
		stop()
	}()

	app := &cli.Command{
		Name:    "claw",
		Usage:   "A downloader and image collector from various sources",
		Version: Version,
		Commands: []*cli.Command{
			internal.ServerCommand(),
		},
		Before: internal.Before,
		After:  internal.After,
		Flags: []cli.Flag{
			internal.ConfigFlag(),
		},
	}

	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
