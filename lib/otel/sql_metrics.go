package otel

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/j2gg0s/otsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var dbmeter = otel.Meter("github.com/tigorlazuardi/claw/lib/otel")

type DBClientDurationMetricHook struct {
	Address string
}

type dbMetricEntry struct {
	start time.Time
}

type dbMetricKey struct{}

func (db *DBClientDurationMetricHook) Before(ctx context.Context, _ *otsql.Event) context.Context {
	return context.WithValue(ctx, dbMetricKey{}, &dbMetricEntry{
		start: time.Now(),
	})
}

var dbClientDurationHistogram, _ = dbmeter.Float64Histogram(
	semconv.DBClientOperationDurationName,
	metric.WithUnit(semconv.DBClientOperationDurationUnit),
	metric.WithDescription("The duration of sqlite database operations used by Claw"),
)

type databaseCallerKey struct{}

func ContextWithDatabaseCaller(ctx context.Context, pc uintptr) context.Context {
	return context.WithValue(ctx, databaseCallerKey{}, pc)
}

func DatabaseCallerFromContext(ctx context.Context) uintptr {
	pc, _ := ctx.Value(databaseCallerKey{}).(uintptr)
	return pc
}

func GetDatabaseOperationName(ctx context.Context) string {
	if pc := DatabaseCallerFromContext(ctx); pc != 0 {
		callers := []uintptr{pc}
		frame, _ := runtime.CallersFrames(callers).Next()
		if frame.Func != nil {
			fnName := frame.Function
			if strings.HasPrefix(fnName, "github.com/tigorlazuardi/claw/") {
				parts := strings.Split(fnName, "/")
				fnName = parts[len(parts)-1]
			}
			return fnName
		}
	}
	return ""
}

func (db *DBClientDurationMetricHook) After(ctx context.Context, evt *otsql.Event) {
	if evt.Method != otsql.MethodQuery && evt.Method != otsql.MethodExec {
		return
	}
	attrs := []attribute.KeyValue{
		semconv.DBSystemSqlite,
		semconv.ServerAddress("file://" + db.Address),
		semconv.DBNamespace(evt.Database),
		attribute.String("method", fmt.Sprintf("sql.conn.%s", evt.Method)),
	}
	if name := GetDatabaseOperationName(ctx); name != "" {
		attrs = append(attrs, semconv.DBOperationName(name))
	}
	if entry, ok := ctx.Value(dbMetricKey{}).(*dbMetricEntry); ok && dbClientDurationHistogram != nil {
		dbClientDurationHistogram.Record(ctx, time.Since(entry.start).Seconds(),
			metric.WithAttributeSet(attribute.NewSet(attrs...)),
		)
	}
}
