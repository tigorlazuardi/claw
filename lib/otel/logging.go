package otel

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

// LoggingProvider wraps the OpenTelemetry LoggerProvider with configuration.
// It embeds *log.LoggerProvider to inherit all its methods and provides
// access to the resolved configuration used during initialization.
type LoggingProvider struct {
	*log.LoggerProvider
	Config LoggingConfig
}

// LoggingOption is a function that configures a LoggingConfig.
// Options are applied in order and can be combined to customize the logging setup.
type LoggingOption func(*LoggingConfig)

// LoggingConfig holds the configuration for the logging provider.
type LoggingConfig struct {
	Exporter log.Exporter
	// Writer is the io.Writer to use for the default exporter.
	// Defaults to os.Stderr. Only used if Exporter is nil.
	Writer io.Writer
	// Resource is the OpenTelemetry resource. If nil, DefaultResource() will be used.
	Resource *resource.Resource
	// Insecure controls whether to use insecure connections for OTLP exporters.
	// Defaults to false (secure connections with proper TLS verification).
	Insecure bool
}

// WithExporter sets a custom log exporter for the logging provider.
// If provided, this exporter will be used instead of the default stderr exporter.
// This allows for custom log destinations such as files, databases, or external services.
func WithExporter(exporter log.Exporter) LoggingOption {
	return func(c *LoggingConfig) {
		c.Exporter = exporter
	}
}

// WithWriter sets the io.Writer for the default log exporter.
// This only affects the default exporter created when no custom exporter is provided.
// Common values are os.Stdout, os.Stderr, or any io.Writer implementation.
// If not provided, defaults to os.Stderr.
func WithWriter(writer io.Writer) LoggingOption {
	return func(c *LoggingConfig) {
		c.Writer = writer
	}
}

// WithResource sets a custom OpenTelemetry resource for the logging provider.
// If provided, this resource will be used instead of the default resource.
// The resource contains service information and other metadata.
func WithResource(resource *resource.Resource) LoggingOption {
	return func(c *LoggingConfig) {
		c.Resource = resource
	}
}

// WithInsecure controls whether to use insecure connections for OTLP exporters.
// When set to true, TLS certificate verification is disabled.
// This should only be used for development or testing environments.
// Defaults to false (secure connections with proper TLS verification).
func WithInsecure(insecure bool) LoggingOption {
	return func(c *LoggingConfig) {
		c.Insecure = insecure
	}
}

// NewLoggingProvider creates a new LoggingProvider with the given options.
// It sets up OpenTelemetry logging with sensible defaults:
//   - Service name: "Claw"
//   - Service version: "0.0.0"
//   - Environment: value of CLAW_ENVIRONMENT env var, or "local" if not set
//   - Exporter: Auto-detects OTLP exporter from environment variables, falls back to stderr
//   - Insecure: false (secure connections with proper TLS verification)
//
// OTLP Exporter Auto-Detection:
//   - Checks OTEL_EXPORTER_OTLP_LOGS_ENDPOINT first, then OTEL_EXPORTER_OTLP_ENDPOINT
//   - Port 4317: Uses gRPC exporter
//   - Port 4318: Uses HTTP exporter
//   - TLS behavior is controlled by the WithInsecure option:
//   - WithInsecure(false): Uses secure connections with proper TLS verification (default)
//   - WithInsecure(true): Disables TLS verification for development/testing
//   - If no OTLP endpoint is configured, falls back to stderr exporter
//
// The created provider is automatically set as the global logger provider.
//
// Example usage:
//
//	// Use defaults
//	provider, err := NewLoggingProvider()
//
//	// Customize service info
//	provider, err := NewLoggingProvider(
//		WithServiceName("my-service"),
//		WithServiceVersion("v2.0.0"),
//		WithEnvironment("production"),
//	)
//
//	// Use custom writer
//	provider, err := NewLoggingProvider(
//		WithWriter(os.Stdout),
//	)
//
//	// Use custom exporter
//	provider, err := NewLoggingProvider(
//		WithExporter(myCustomExporter),
//	)
//
//	// OTLP auto-detection examples:
//	// export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317  # Uses gRPC
//	// export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318  # Uses HTTP
//	// export OTEL_EXPORTER_OTLP_LOGS_ENDPOINT=https://otel.example.com:4317  # gRPC with TLS
//	provider, err := NewLoggingProvider()  // Will auto-detect from environment
func NewLoggingProvider(opts ...LoggingOption) (*LoggingProvider, error) {
	config := &LoggingConfig{
		Writer: os.Stderr,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.Exporter == nil {
		// First try to detect OTLP exporter from environment variables
		otlpExporter, err := detectOTLPLogExporter(config.Insecure)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
		}

		if otlpExporter != nil {
			config.Exporter = otlpExporter
		} else {
			// Fall back to stderr exporter if no OTLP endpoint is configured
			exporter, err := stdoutlog.New(
				stdoutlog.WithWriter(config.Writer),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create log exporter: %w", err)
			}
			config.Exporter = exporter
		}
	}

	if config.Resource == nil {
		res, err := DefaultResource()
		if err != nil {
			return nil, fmt.Errorf("failed to create default resource: %w", err)
		}
		config.Resource = res
	}

	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(config.Exporter)),
		log.WithResource(config.Resource),
	)

	return &LoggingProvider{
		LoggerProvider: provider,
		Config:         *config,
	}, nil
}

// detectOTLPLogExporter detects OTLP endpoint from environment variables and creates appropriate exporter.
// It checks OTEL_EXPORTER_OTLP_LOGS_ENDPOINT first, then falls back to OTEL_EXPORTER_OTLP_ENDPOINT.
// Returns nil if no OTLP endpoint is configured.
func detectOTLPLogExporter(insecure bool) (log.Exporter, error) {
	// Check for logs-specific endpoint first
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT")
	if endpoint == "" {
		// Fall back to general OTLP endpoint
		endpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}

	if endpoint == "" {
		return nil, nil // No OTLP endpoint configured
	}

	// Parse the URL to determine the port
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid OTLP endpoint URL: %w", err)
	}

	// Determine protocol based on port
	var port int
	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid port in OTLP endpoint: %w", err)
		}
	}

	// Create exporter based on port
	switch port {
	case 4317:
		// gRPC exporter
		return createGRPCLogExporter(endpoint, insecure)
	case 4318:
		// HTTP exporter
		return createHTTPLogExporter(endpoint, insecure)
	default:
		// Default to gRPC for unknown ports
		return createGRPCLogExporter(endpoint, insecure)
	}
}

// createGRPCLogExporter creates an OTLP gRPC log exporter.
func createGRPCLogExporter(endpoint string, insecure bool) (log.Exporter, error) {
	// Configure gRPC connection options
	opts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
	}

	if insecure {
		// Use insecure connection when explicitly requested
		opts = append(opts, otlploggrpc.WithTLSCredentials(grpcinsecure.NewCredentials()))
	} else {
		// Use secure connection with proper TLS verification
		opts = append(opts, otlploggrpc.WithTLSCredentials(
			credentials.NewTLS(&tls.Config{}),
		))
	}

	return otlploggrpc.New(context.Background(), opts...)
}

// createHTTPLogExporter creates an OTLP HTTP log exporter.
func createHTTPLogExporter(endpoint string, insecure bool) (log.Exporter, error) {
	// Configure HTTP client options
	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(endpoint),
	}

	if insecure {
		// Use insecure connection when explicitly requested
		opts = append(opts, otlploghttp.WithInsecure())
	} else {
		// Use secure connection with proper TLS verification
		opts = append(opts, otlploghttp.WithTLSClientConfig(&tls.Config{}))
	}

	return otlploghttp.New(context.Background(), opts...)
}

// getEnvironment returns the deployment environment from the CLAW_ENVIRONMENT
// environment variable, or "local" if not set.
func getEnvironment() string {
	if env := os.Getenv("CLAW_ENVIRONMENT"); env != "" {
		return env
	}
	return "local"
}
