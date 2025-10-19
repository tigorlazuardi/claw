package otel

import (
	"context"
	"os"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

func GetTraceEndpoint() string {
	traceEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if traceEndpoint == "" {
		traceEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	return traceEndpoint
}

func SetupTracing(ctx context.Context) error {
	if endpoint := GetTraceEndpoint(); endpoint != "" {
		provider, err := CreateTraceProvider(ctx)
		if err != nil {
			return err
		}
		otel.SetTracerProvider(provider)
		return nil
	}
	otel.SetTracerProvider(CreateNoopTraceProvider())
	return nil
}

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
	Shutdowns = append(Shutdowns, batcher.ForceFlush, batcher.Shutdown)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(Resource),
	), nil
}

// CreateNoopTraceProvider returns a a tracer provider that does not export any traces,
// but can be used to create spans. Useful for connecting logs with traces even if traces
// are not used by the user.
func CreateNoopTraceProvider() trace.TracerProvider {
	return sdktrace.NewTracerProvider()
}

var tracer = otel.Tracer("github.com/tigorlazuardi/claw/lib/otel")

func StartSpan(ctx context.Context) (context.Context, trace.Span) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	frame, _ := runtime.CallersFrames(pc).Next()
	fnName := frame.Function
	if strings.HasPrefix(frame.Function, "github.com/tigorlazuardi/claw/") {
		parts := strings.Split(frame.Function, "/")
		fnName = parts[len(parts)-1]
	}
	opts := trace.WithAttributes(
		semconv.CodeFunctionName(frame.Function),
		semconv.CodeFilePath(frame.File),
		semconv.CodeLineNumber(frame.Line),
	)
	ctx, span := tracer.Start(ctx, fnName, opts)
	return ctx, span
}
