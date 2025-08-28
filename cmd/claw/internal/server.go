package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/tigorlazuardi/claw/lib/claw"
	"github.com/tigorlazuardi/claw/lib/server"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	_ "modernc.org/sqlite"
)

// ServerCommand creates the server CLI command
func ServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the claw server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "Server address to listen on",
				Value:   ":8080",
				Sources: cli.EnvVars("CLAW_SERVER_ADDR"),
			},
			&cli.StringFlag{
				Name:    "db-path",
				Usage:   "Path to SQLite database file",
				Value:   "./claw.db",
				Sources: cli.EnvVars("CLAW_DB_PATH"),
			},
		},
		Action: runServer,
	}
}

// runServer starts the HTTP server with ConnectRPC handlers
func runServer(ctx context.Context, cmd *cli.Command) error {
	addr := cmd.String("addr")
	dbPath := cmd.String("db-path")

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Initialize the claw service
	clawService := claw.New(db)

	// Create handlers
	sourceHandler := server.NewSourceHandler(clawService)
	deviceHandler := server.NewDeviceHandler(clawService)
	imageHandler := server.NewImageHandler(clawService)
	tagHandler := server.NewTagHandler(clawService)

	// Create HTTP mux and register ConnectRPC handlers
	mux := http.NewServeMux()
	
	// Register all service handlers
	sourcePath, sourceHandlerHTTP := clawv1connect.NewSourceServiceHandler(sourceHandler)
	mux.Handle(sourcePath, sourceHandlerHTTP)
	
	devicePath, deviceHandlerHTTP := clawv1connect.NewDeviceServiceHandler(deviceHandler)
	mux.Handle(devicePath, deviceHandlerHTTP)
	
	imagePath, imageHandlerHTTP := clawv1connect.NewImageServiceHandler(imageHandler)
	mux.Handle(imagePath, imageHandlerHTTP)
	
	tagPath, tagHandlerHTTP := clawv1connect.NewTagServiceHandler(tagHandler)
	mux.Handle(tagPath, tagHandlerHTTP)

	// Create HTTP server with h2c support for HTTP/2 over cleartext
	httpServer := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server", "addr", addr)
		errChan <- httpServer.ListenAndServe()
	}()

	// Wait for either server error or shutdown signal
	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			return fmt.Errorf("server failed: %w", err)
		}
	case <-ctx.Done():
		slog.Info("Received shutdown signal, shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Graceful shutdown
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
			return err
		}
		slog.Info("Server shutdown complete")
	}
	return nil
}

