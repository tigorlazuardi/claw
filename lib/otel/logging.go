package otel

import (
	"fmt"
	"io"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
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
	// ServiceName is the name of the service. Defaults to "Claw".
	ServiceName string
	// ServiceVersion is the version of the service. Defaults to "0.0.0".
	ServiceVersion string
	// Environment is the deployment environment. Defaults to the value of
	// CLAW_ENVIRONMENT environment variable, or "local" if not set.
	Environment string
	// Exporter is the log exporter to use. If nil, a default stderr exporter
	// using stdoutlog.New(stdoutlog.WithWriter(Writer)) will be created.
	Exporter log.Exporter
	// Writer is the io.Writer to use for the default exporter.
	// Defaults to os.Stderr. Only used if Exporter is nil.
	Writer io.Writer
	// Resource is the OpenTelemetry resource. If nil, DefaultResource() will be used.
	Resource *resource.Resource
}

// WithServiceName sets the service name for the logging provider.
// This will be used in the OpenTelemetry resource attributes.
// If not provided, defaults to "Claw".
func WithServiceName(name string) LoggingOption {
	return func(c *LoggingConfig) {
		c.ServiceName = name
	}
}

// WithServiceVersion sets the service version for the logging provider.
// This will be used in the OpenTelemetry resource attributes.
// If not provided, defaults to "v1".
func WithServiceVersion(version string) LoggingOption {
	return func(c *LoggingConfig) {
		c.ServiceVersion = version
	}
}

// WithEnvironment sets the deployment environment for the logging provider.
// This will be used in the OpenTelemetry resource attributes.
// If not provided, defaults to the value of CLAW_ENVIRONMENT environment variable,
// or "local" if the environment variable is not set.
func WithEnvironment(env string) LoggingOption {
	return func(c *LoggingConfig) {
		c.Environment = env
	}
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

// NewLoggingProvider creates a new LoggingProvider with the given options.
// It sets up OpenTelemetry logging with sensible defaults:
//   - Service name: "Claw"
//   - Service version: "v1"
//   - Environment: value of CLAW_ENVIRONMENT env var, or "local" if not set
//   - Exporter: stderr exporter using stdoutlog.New(stdoutlog.WithWriter(os.Stderr))
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
func NewLoggingProvider(opts ...LoggingOption) (*LoggingProvider, error) {
	config := &LoggingConfig{
		ServiceName:    "Claw",
		ServiceVersion: "0.0.0",
		Environment:    getEnvironment(),
		Writer:         os.Stderr,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.Exporter == nil {
		exporter, err := stdoutlog.New(
			stdoutlog.WithWriter(config.Writer),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create log exporter: %w", err)
		}
		config.Exporter = exporter
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

// getEnvironment returns the deployment environment from the CLAW_ENVIRONMENT
// environment variable, or "local" if not set.
func getEnvironment() string {
	if env := os.Getenv("CLAW_ENVIRONMENT"); env != "" {
		return env
	}
	return "local"
}
