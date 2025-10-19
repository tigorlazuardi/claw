package otel

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var PrometheusExporter *prometheus.Exporter

func GetMetricsEndpoint() string {
	metricsEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT")
	if metricsEndpoint == "" {
		metricsEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	return metricsEndpoint
}

func SetupMetrics(ctx context.Context) error {
	opts := []sdkmetric.Option{
		sdkmetric.WithResource(Resource),
	}
	if endpoint := GetMetricsEndpoint(); endpoint != "" {
		reader, err := CreateMetricsReader(ctx)
		if err != nil {
			otel.Handle(err)
		} else {
			opts = append(opts, metric.WithReader(reader))
		}
	}

	if b, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("CLAW_PROMETHEUS_ENABLE"))); b {
		prom, err := prometheus.New(
			prometheus.WithResourceAsConstantLabels(func(kv attribute.KeyValue) bool { return true }),
		)
		if err != nil {
			otel.Handle(err)
		} else {
			Shutdowns = append(Shutdowns, prom.Shutdown)
			opts = append(opts, metric.WithReader(prom))
			PrometheusExporter = prom
		}
	}

	provider := metric.NewMeterProvider(opts...)
	otel.SetMeterProvider(provider)
	return nil
}

func CreateMetricsReader(ctx context.Context) (sdkmetric.Reader, error) {
	proto := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL")
	if proto == "" {
		proto = os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	}

	var (
		exporter sdkmetric.Exporter
		err      error
	)
	if proto == "grpc" {
		exporter, err = otlpmetricgrpc.New(ctx)
	} else {
		exporter, err = otlpmetrichttp.New(ctx)
	}
	if err != nil {
		return nil, err
	}
	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Second*5))
	Shutdowns = append(Shutdowns, reader.ForceFlush, reader.Shutdown)
	return reader, nil
}
