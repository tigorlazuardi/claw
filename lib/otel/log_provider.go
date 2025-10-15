package otel

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func CreateLogProvider(ctx context.Context) (log.LoggerProvider, error) {
	proto := os.Getenv("OTEL_EXPORTER_OTLP_LOGS_PROTOCOL")
	if proto == "" {
		proto = os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	}
	var (
		exporter sdklog.Exporter
		err      error
	)
	if proto == "grpc" {
		exporter, err = otlploggrpc.New(ctx)
	} else {
		exporter, err = otlploghttp.New(ctx)
	}
	if err != nil {
		return nil, err
	}
	batcher := sdklog.NewBatchProcessor(exporter)
	return sdklog.NewLoggerProvider(sdklog.WithProcessor(batcher)), nil
}
