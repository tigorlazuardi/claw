package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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
		Name:   "server",
		Usage:  "Start the claw server",
		Action: runServer,
	}
}

// runServer starts the HTTP server with ConnectRPC handlers
func runServer(ctx context.Context, cmd *cli.Command) error {
	// Open database connection
	db, err := sql.Open("sqlite", cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	abs, _ := filepath.Abs(cfg.Database.Path)
	slog.Info("Database connected", "path", abs)

	// Initialize the claw service
	clawService := claw.New(db, cfg.Claw)

	cfg.OnChange = clawService.RereadConfig
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGUSR1)
		defer signal.Stop(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				slog.Info("received SIGUSR1, reloading configuration")
				clawService.RereadConfig()
			}
		}
	}()

	// Create handlers
	sourceHandler := server.NewSourceHandler(clawService)
	deviceHandler := server.NewDeviceHandler(clawService)
	imageHandler := server.NewImageHandler(clawService)
	tagHandler := server.NewTagHandler(clawService)
	jobHandler := server.NewJobHandler(clawService)

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

	jobPath, jobHandlerHTTP := clawv1connect.NewJobServiceHandler(jobHandler)
	mux.Handle(jobPath, jobHandlerHTTP)

	listener := cfg.Server.Host
	if listener == nil {
		ln, err := net.Listen("tcp", ":8000")
		if err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}
		listener = &NetListener{Listener: ln}
	}

	// Create HTTP server with h2c support for HTTP/2 over cleartext
	httpServer := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		if listener.Addr().Network() == "tcp" {
			port := listener.Addr().(*net.TCPAddr).Port
			outgoingAddr, err := net.Dial("udp", "1.1.1.1:53")
			if err != nil {
				errChan <- fmt.Errorf("failed to determine outgoing address: %w", err)
				return
			}
			publicIp := outgoingAddr.LocalAddr().(*net.UDPAddr).IP.String()
			slog.Info(fmt.Sprintf("Server outgoing address: http://%s:%d", publicIp, port))
		}
		slog.Info("Server is listening", "address", listener.Addr().String())
		errChan <- httpServer.Serve(listener)
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
