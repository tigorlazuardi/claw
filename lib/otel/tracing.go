package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// TracingProvider wraps the OpenTelemetry TracerProvider with configuration.
// It embeds *trace.TracerProvider to inherit all its methods and provides
// access to the resolved configuration used during initialization.
type TracingProvider struct {
	*trace.TracerProvider
	Config TracingConfig
}

// TracingOption is a function that configures a TracingConfig.
// Options are applied in order and can be combined to customize the tracing setup.
type TracingOption func(*TracingConfig)

// TracingConfig holds the configuration for the tracing provider.
// All fields are public to allow direct access to the resolved configuration.
type TracingConfig struct {
	// ServiceName is the name of the service. Defaults to "Claw".
	ServiceName string
	// ServiceVersion is the version of the service. Defaults to "0.0.0".
	ServiceVersion string
	// Environment is the deployment environment. Defaults to the value of
	// CLAW_ENVIRONMENT environment variable, or "local" if not set.
	Environment string
	// Exporter is the trace exporter to use. If nil, a NOOP exporter will be used.
	Exporter trace.SpanExporter
	// Resource is the OpenTelemetry resource. If nil, DefaultResource() will be used.
	Resource *resource.Resource
	// Sampler is the trace sampler. If nil, trace.AlwaysSample() will be used.
	Sampler trace.Sampler
}

// WithTracingServiceName sets the service name for the tracing provider.
// This will be used in the OpenTelemetry resource attributes.
// If not provided, defaults to "Claw".
func WithTracingServiceName(name string) TracingOption {
	return func(c *TracingConfig) {
		c.ServiceName = name
	}
}

// WithTracingServiceVersion sets the service version for the tracing provider.
// This will be used in the OpenTelemetry resource attributes.
// If not provided, defaults to "0.0.0".
func WithTracingServiceVersion(version string) TracingOption {
	return func(c *TracingConfig) {
		c.ServiceVersion = version
	}
}

// WithTracingEnvironment sets the deployment environment for the tracing provider.
// This will be used in the OpenTelemetry resource attributes.
// If not provided, defaults to the value of CLAW_ENVIRONMENT environment variable,
// or "local" if the environment variable is not set.
func WithTracingEnvironment(env string) TracingOption {
	return func(c *TracingConfig) {
		c.Environment = env
	}
}

// WithTracingExporter sets a custom trace exporter for the tracing provider.
// If provided, this exporter will be used instead of the default NOOP exporter.
// This allows for custom trace destinations such as OTLP endpoints, stdout, files, etc.
func WithTracingExporter(exporter trace.SpanExporter) TracingOption {
	return func(c *TracingConfig) {
		c.Exporter = exporter
	}
}

// WithTracingResource sets a custom OpenTelemetry resource for the tracing provider.
// If provided, this resource will be used instead of the default resource.
// The resource contains service information and other metadata.
func WithTracingResource(resource *resource.Resource) TracingOption {
	return func(c *TracingConfig) {
		c.Resource = resource
	}
}

// WithTracingSampler sets a custom trace sampler for the tracing provider.
// If provided, this sampler will be used instead of the default AlwaysSample sampler.
// Common samplers include trace.AlwaysSample(), trace.NeverSample(), and trace.TraceIDRatioBased().
func WithTracingSampler(sampler trace.Sampler) TracingOption {
	return func(c *TracingConfig) {
		c.Sampler = sampler
	}
}

// NewTracingProvider creates a new TracingProvider with the given options.
// It sets up OpenTelemetry tracing with sensible defaults:
//   - Service name: "Claw"
//   - Service version: "0.0.0"
//   - Environment: value of CLAW_ENVIRONMENT env var, or "local" if not set
//   - Exporter: NOOP exporter (no traces are exported)
//   - Resource: DefaultResource()
//   - Sampler: AlwaysSample()
//
// The created provider is automatically set as the global tracer provider.
//
// Example usage:
//
//	// Use defaults (NOOP exporter)
//	provider, err := NewTracingProvider()
//
//	// Customize with OTLP exporter
//	provider, err := NewTracingProvider(
//		WithExporter(otlpExporter),
//		WithServiceName("my-service"),
//	)
func NewTracingProvider(opts ...TracingOption) (*TracingProvider, error) {
	config := &TracingConfig{
		ServiceName:    "Claw",
		ServiceVersion: "0.0.0",
		Environment:    getEnvironment(),
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.Exporter == nil {
		config.Exporter = tracetest.NewNoopExporter()
	}

	if config.Resource == nil {
		res, err := DefaultResource()
		if err != nil {
			return nil, fmt.Errorf("failed to create default resource: %w", err)
		}
		config.Resource = res
	}

	if config.Sampler == nil {
		config.Sampler = trace.AlwaysSample()
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(config.Exporter),
		trace.WithResource(config.Resource),
		trace.WithSampler(config.Sampler),
	)

	otel.SetTracerProvider(provider)

	return &TracingProvider{
		TracerProvider: provider,
		Config:         *config,
	}, nil
}

// Shutdown gracefully shuts down the tracing provider.
// It should be called when the application is terminating to ensure
// all trace spans are properly flushed and exported.
func (tp *TracingProvider) Shutdown(ctx context.Context) error {
	return tp.TracerProvider.Shutdown(ctx)
}