package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type JSONHandler struct {
	*slog.JSONHandler
}

func (js *JSONHandler) Handle(ctx context.Context, rec slog.Record) error {
	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		rec.AddAttrs(
			slog.String("trace_id", span.TraceID().String()),
			slog.String("span_id", span.SpanID().String()),
		)
	}
	return js.JSONHandler.Handle(ctx, rec)
}

func (js *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JSONHandler{
		JSONHandler: js.JSONHandler.WithAttrs(attrs).(*slog.JSONHandler),
	}
}

func (js *JSONHandler) WithGroup(name string) slog.Handler {
	return &JSONHandler{
		JSONHandler: js.JSONHandler.WithGroup(name).(*slog.JSONHandler),
	}
}
