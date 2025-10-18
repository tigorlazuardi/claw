package logger

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"
	"unsafe"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ slog.Handler = (*OtelHandler)(nil)

type OtelHandler struct {
	*otelslog.Handler
}

func transformOtelAttr(groups []string, attr slog.Attr) slog.Attr {
	if _, custom := attr.Value.Any().(slog.LogValuer); custom {
		return attr
	}
	if attr.Value.Kind() == slog.KindDuration {
		attr.Value = slog.StringValue(attr.Value.Duration().String())
		return attr
	}
	if attr.Value.Kind() == slog.KindGroup {
		groups = append(groups, attr.Key)
		attrs := attr.Value.Group()
		for i := range attrs {
			attrs[i] = transformOtelAttr(groups, attrs[i])
		}
		return attr
	}
	switch val := attr.Value.Any().(type) {
	case error:
		attr.Value = slog.GroupValue(
			slog.Group("exception",
				slog.String("message", val.Error()),
				slog.String("type", fmt.Sprintf("%T", val)),
			),
		)
		return attr
	case []byte:
		if utf8.Valid(val) {
			attr.Value = slog.StringValue(unsafe.String(unsafe.SliceData(val), len(val)))
		} else {
			attr.Value = slog.StringValue("!BINARY:" + base64.StdEncoding.EncodeToString(val))
		}
		return attr
	case proto.Message:
		b, err := protojson.Marshal(val)
		if err != nil {
			attr.Value = slog.StringValue("!ERROR:" + err.Error())
		} else {
			attr.Value = slog.StringValue(unsafe.String(unsafe.SliceData(b), len(b)))
		}
		return attr
	default:
		b, err := json.Marshal(val)
		if err == nil {
			attr.Value = slog.StringValue(unsafe.String(unsafe.SliceData(b), len(b)))
		} else {
			attr.Value = slog.StringValue("!ERROR:" + err.Error())
		}
		return attr
	}
}

var wd, _ = os.Getwd()

func (ot *OtelHandler) Handle(ctx context.Context, rec slog.Record) error {
	r := slog.NewRecord(rec.Time, rec.Level, rec.Message, rec.PC)
	if r.PC != 0 {
		if frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next(); frame.Func != nil {
			fnParts := strings.Split(frame.Function, "/")
			file := frame.File
			if wd != "" {
				file = strings.TrimPrefix(frame.File, wd+"/")
			}
			r.AddAttrs(slog.Group("code",
				slog.Group("function", slog.String("name", fnParts[len(fnParts)-1])),
				slog.Group("file", slog.String("path", file)),
				slog.Group("line", slog.Int("number", frame.Line)),
			))
		}
	}
	rec.Attrs(func(a slog.Attr) bool {
		r.AddAttrs(transformOtelAttr([]string{}, a))
		return true
	})
	return ot.Handler.Handle(ctx, r)
}

func (ot *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &OtelHandler{
		Handler: ot.Handler.WithAttrs(attrs).(*otelslog.Handler),
	}
}

func (ot *OtelHandler) WithGroup(name string) slog.Handler {
	return &OtelHandler{
		Handler: ot.Handler.WithGroup(name).(*otelslog.Handler),
	}
}
