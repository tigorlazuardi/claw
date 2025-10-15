package otel

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func CreateTraceProvider(ctx context.Context) (trace.TracerProvider, error) {
	proto := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL")
	if proto == "" {
		proto = os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	}
	var (
		exporter sdktrace.SpanExporter
		err      error
	)
	if proto == "grpc" {
		exporter, err = otlptracegrpc.New(ctx)
	} else {
		exporter, err = otlptracehttp.New(ctx)
	}
	if err != nil {
		return nil, err
	}
	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	return sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(batcher)), nil
}

// CreateNoopTraceProvider returns a no-op TracerProvider.
//
// This will allow tracing spans to be created but they will not be exported anywhere.
//
// Useful for logging using log/slog since they will be included in the log entry if they are printed to stderr, stdout, or a file.
func CreateNoopTraceProvider() trace.TracerProvider {
	return noop.TracerProvider{}
}
