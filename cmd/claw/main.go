package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tigorlazuardi/claw/cmd/claw/internal"
	"github.com/urfave/cli/v3"
)

func main() {
	// Create context that cancels on interrupt/termination signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	app := &cli.Command{
		Name:    "claw",
		Usage:   "A downloader and image collector from various sources",
		Version: "1.0.0",
		Commands: []*cli.Command{
			internal.ServerCommand(),
		},
	}

	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

