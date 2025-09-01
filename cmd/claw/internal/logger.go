package internal

import (
	"context"
	"log/slog"
)

// LogHandler is a slog.LogHandler that groups all non-standard attributes into a "details" group
// so it can be consumed better by log aggregators.
type LogHandler struct {
	slog.Handler
}

func (lo LogHandler) Handle(ctx context.Context, record slog.Record) error {
	rec := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	details := make([]slog.Attr, 0, 4)
	record.Attrs(func(a slog.Attr) bool {
		// Slog Keys are already filtered by the standard library
		details = append(details, a)
		return true
	})
	rec.AddAttrs(slog.Attr{Key: "details", Value: slog.GroupValue(details...)})
	return lo.Handler.Handle(ctx, rec)
}

func (lo LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LogHandler{
		Handler: lo.Handler.WithAttrs(attrs),
	}
}

func (lo LogHandler) WithGroup(name string) slog.Handler {
	return LogHandler{
		Handler: lo.Handler.WithGroup(name),
	}
}
