package logger

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/j2gg0s/otsql"
)

type skipLogKey struct{}

func ContextWithSkipLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipLogKey{}, struct{}{})
}

type LoggerHook struct {
	Logger *slog.Logger
}

func (lo LoggerHook) Before(ctx context.Context, _ *otsql.Event) context.Context {
	return ctx
}

func (lo LoggerHook) After(ctx context.Context, evt *otsql.Event) {
	if ctx.Value(skipLogKey{}) != nil {
		return
	}
	if evt.Method != otsql.MethodQuery && evt.Method != otsql.MethodExec {
		return
	}
	lvl := slog.LevelInfo
	end := time.Now()
	dur := end.Sub(evt.BeginAt)
	if evt.Err != nil {
		if errors.Is(evt.Err, driver.ErrSkip) {
			return
		}
		lvl = slog.LevelError
	}
	if dur > time.Second*3 {
		lvl = slog.LevelWarn
	}
	handler := lo.Logger.Handler()
	if !handler.Enabled(ctx, lvl) {
		return
	}
	attrs := []slog.Attr{}
	if evt.Args != nil {
		attrs = append(attrs, slog.Any("args", argsLogValue{args: evt.Args}))
	}
	if evt.Err != nil {
		attrs = append(attrs, slog.Any("error", errorLogValue{err: evt.Err}))
	}
	attrs = append(attrs, slog.Time("begin_at", evt.BeginAt))
	attrs = append(attrs, slog.Time("end_at", end))
	attrs = append(attrs, slog.Duration("duration", dur))
	if evt.Method != "" {
		attrs = append(attrs, slog.String("method", string(evt.Method)))
	}
	if evt.Database != "" {
		attrs = append(attrs, slog.String("database", evt.Database))
	}
	if evt.Conn != "" {
		attrs = append(attrs, slog.String("conn", evt.Conn))
	}
	rec := slog.NewRecord(end, lvl, strings.TrimSpace(evt.Query), 0)
	rec.AddAttrs(attrs...)
	handler.Handle(ctx, rec)
}

type errorLogValue struct {
	err error
}

func (elv errorLogValue) LogValue() slog.Value {
	if elv.err == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.GroupAttrs("exception",
			slog.String("message", elv.err.Error()),
			slog.String("type", fmt.Sprintf("%T", elv.err)),
		),
	)
}

type argsLogValue struct {
	args any
}

func (ar argsLogValue) LogValue() slog.Value {
	switch v := ar.args.(type) {
	case []driver.NamedValue:
		return namedValueLogger(v).LogValue()
	case []driver.Value:
		return valueLogger(v).LogValue()
	default:
		return slog.AnyValue(ar.args)
	}
}

type namedValueLogger []driver.NamedValue

func (va namedValueLogger) LogValue() slog.Value {
	attrs := make([]slog.Attr, len(va))
	for i, v := range va {
		if v.Name != "" {
			attrs[i] = slog.Any(v.Name, v.Value)
		} else {
			attrs[i] = slog.Any(strconv.Itoa(v.Ordinal), v.Value)
		}
	}
	return slog.GroupValue(attrs...)
}

type valueLogger []driver.Value

func (va valueLogger) LogValue() slog.Value {
	attrs := make([]slog.Attr, len(va))
	for i, v := range va {
		attrs[i] = slog.Any(strconv.Itoa(i+1), v)
	}
	return slog.GroupValue(attrs...)
}
