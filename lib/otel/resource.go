package otel

import (
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// DefaultResource creates a default OpenTelemetry resource with sensible defaults.
// It sets:
//   - Service name: "Claw"
//   - Service version: "0.0.0"
//   - Environment: value of CLAW_ENVIRONMENT env var, or "local" if not set
//
// The resource is merged with the default resource provided by the SDK.
func DefaultResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("Claw"),
			semconv.ServiceVersion("0.0.0"),
			semconv.DeploymentEnvironment(getEnvironment()),
		),
	)
}

