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

	"connectrpc.com/connect"
	"github.com/networkteam/go-sqllogger"
	"github.com/tigorlazuardi/claw/lib/claw"
	"github.com/tigorlazuardi/claw/lib/dblogger"
	"github.com/tigorlazuardi/claw/lib/server"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
	"github.com/tigorlazuardi/claw/migrations"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	sqlite "github.com/ncruces/go-sqlite3/driver"

	_ "github.com/ncruces/go-sqlite3/embed"
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
	conn, err := (&sqlite.SQLite{}).OpenConnector(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	ldb := sqllogger.LoggingConnector(dblogger.DBLogger{
		Logger: slog.Default(),
	}, conn)
	db := sql.OpenDB(ldb)
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

	// Create logging interceptor
	loggingInterceptor := server.LoggingInterceptor(slog.Default())

	// Register all service handlers with interceptor
	sourcePath, sourceHandlerHTTP := clawv1connect.NewSourceServiceHandler(sourceHandler,
		connect.WithInterceptors(loggingInterceptor))
	mux.Handle(sourcePath, sourceHandlerHTTP)

	devicePath, deviceHandlerHTTP := clawv1connect.NewDeviceServiceHandler(deviceHandler,
		connect.WithInterceptors(loggingInterceptor))
	mux.Handle(devicePath, deviceHandlerHTTP)

	imagePath, imageHandlerHTTP := clawv1connect.NewImageServiceHandler(imageHandler,
		connect.WithInterceptors(loggingInterceptor))
	mux.Handle(imagePath, imageHandlerHTTP)

	tagPath, tagHandlerHTTP := clawv1connect.NewTagServiceHandler(tagHandler,
		connect.WithInterceptors(loggingInterceptor))
	mux.Handle(tagPath, tagHandlerHTTP)

	jobPath, jobHandlerHTTP := clawv1connect.NewJobServiceHandler(jobHandler,
		connect.WithInterceptors(loggingInterceptor))
	mux.Handle(jobPath, jobHandlerHTTP)

	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}))

	webuiFragment, distFS, err := CreateViteFragment()
	if err != nil {
		return fmt.Errorf("failed to create Vite fragment: %w", err)
	}

	// Create WebUI handler with logging middleware
	webuiHandler := CreateWebuiHandler(WebUIConfig{
		Fragment: webuiFragment,
		DistFS:   distFS,
		Logger:   slog.Default(),
	})

	// Wrap WebUI handler with HTTP logging middleware
	httpLoggingMiddleware := server.HTTPLoggingMiddleware(slog.Default())
	mux.Handle("/", httpLoggingMiddleware(webuiHandler))

	listener := cfg.Server.Host
	if listener == nil {
		ln, err := net.Listen("tcp", ":8000")
		if err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}
		listener = &NetListener{Listener: ln}
	}

	corsMiddleware := corsDevMidddlware(cfg.Server.WebUI.DevMode)

	// Create HTTP server with h2c support for HTTP/2 over cleartext
	httpServer := &http.Server{
		Handler: corsMiddleware(stripPrefixHandler(h2c.NewHandler(mux, &http2.Server{}))),
	}

	if err := migrations.Migrate(ctx, db); err != nil {
		return err
	}

	errChan := make(chan error, 1)
	schedulerExit := make(chan struct{}, 1)
	go func() {
		clawService.StartSchedculer(ctx)
		schedulerExit <- struct{}{}
	}()

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
			slog.Info("server is listening",
				"public", fmt.Sprintf("http://%s:%d", publicIp, port),
				"local", fmt.Sprintf("http://%s", listener.Addr().String()),
			)
		} else { // unix socket
			slog.Info("server is listening", "socket", listener.Addr().String())
		}
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
		<-schedulerExit
		slog.Info("Server shutdown complete")
	}
	return nil
}

func stripPrefixHandler(handler http.Handler) http.Handler {
	if cfg.Server.BaseURL == "" || cfg.Server.BaseURL == "/" {
		return handler
	}
	slog.Info("Serving under base URL", "base_url", cfg.Server.BaseURL)
	return http.StripPrefix(cfg.Server.BaseURL, handler)
}
