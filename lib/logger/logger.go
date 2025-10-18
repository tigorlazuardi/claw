package logger

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"unicode/utf8"
	"unsafe"

	clawotel "github.com/tigorlazuardi/claw/lib/otel"
	"github.com/tigorlazuardi/prettylog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
)

func Setup(ctx context.Context) error {
	if clawotel.IsLogEndpointSet() {
		provider, err := clawotel.CreateLogProvider(ctx)
		if err != nil {
			return err
		}
		global.SetLoggerProvider(provider)
		handler := otelslog.NewHandler("github.com/tigorlazuardi/claw")
		slog.SetDefault(slog.New(&OtelHandler{Handler: handler}))
		return nil
	}
	if prettylog.CanColor(os.Stderr) {
		prettyHandler := prettylog.New(
			prettylog.WithPackageName("github.com/tigorlazuardi/claw"),
			prettylog.AddWritersBefore(prettylog.DefaultPrettyJSONWriter, prettylog.NewCommonWriter(func(info prettylog.RecordData) string {
				span := trace.SpanContextFromContext(info.Context)
				if span.HasTraceID() {
					return span.TraceID().String()
				}
				return ""
			}).WithStaticKey("Trace ID")),
			prettylog.AddWritersBefore(prettylog.DefaultPrettyJSONWriter, prettylog.NewCommonWriter(func(info prettylog.RecordData) string {
				span := trace.SpanContextFromContext(info.Context)
				if span.HasSpanID() {
					return span.SpanID().String()
				}
				return ""
			}).WithStaticKey("Span ID")),
		)
		slog.SetDefault(slog.New(prettyHandler))
		return nil
	}
	logger := slog.New(&JSONHandler{slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelInfo,
		ReplaceAttr: replaceAttr,
	})})
	slog.SetDefault(logger)
	return nil
}

var cwd, _ = os.Getwd()

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if _, logValuer := a.Value.Any().(slog.LogValuer); logValuer {
		return a
	}
	if a.Value.Kind() == slog.KindDuration {
		a.Value = slog.StringValue(a.Value.Duration().String())
	}
	if b, ok := a.Value.Any().([]byte); ok {
		if utf8.Valid(b) {
			a.Value = slog.StringValue(unsafe.String(unsafe.SliceData(b), len(b)))
		} else {
			a.Value = slog.StringValue("!BINARY:" + base64.StdEncoding.EncodeToString(b))
		}
		return a
	}
	if source, ok := a.Value.Any().(*slog.Source); ok {
		source.File = strings.TrimPrefix(source.File, cwd+string(os.PathSeparator))
		source.Function = strings.TrimPrefix(source.Function, "github.com/tigorlazuardi/claw/")
	}
	if m, ok := a.Value.Any().(proto.Message); ok {
		a.Value = transformProtoToLog(m)
	}
	return a
}

func transformProtoToLog(msg proto.Message) slog.Value {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return slog.StringValue("!ERROR:" + err.Error())
	}
	return slog.AnyValue(json.RawMessage(b))
}
