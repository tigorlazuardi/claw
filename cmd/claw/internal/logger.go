package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"

	"github.com/mattn/go-isatty"
	"github.com/tidwall/pretty"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var bufPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

var poolMemoryCounter atomic.Int64

const maxLogBufferMemory = 8 * 1024 * 1024 // 8MB

// LogHandler is a slog.LogHandler that groups all non-standard attributes into a "details" group
// so it can be consumed better by log aggregators.
type LogHandler struct {
	attrs []slog.Attr
	group []string
	level slog.Level
	mu    sync.Mutex
}

var canPretty = isatty.IsTerminal(os.Stderr.Fd())

func (lo *LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= lo.level
}

func (lo *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	rec := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	details := make([]slog.Attr, 0, 4)
	record.Attrs(func(a slog.Attr) bool {
		// Slog Keys are already filtered by the standard library
		details = append(details, a)
		return true
	})
	rec.AddAttrs(slog.Attr{Key: "details", Value: slog.GroupValue(details...)})
	buf := bufPool.Get().(*bytes.Buffer)
	poolMemoryCounter.Add(-int64(buf.Cap()))
	defer func() {
		if buf.Cap() > 8*1024 {
			// Discard buffers larger than 8KB, so large buffers
			// will not be kept in the pool.
			return
		}
		// only put back to pool if we are under the max memory limit
		// otherwise let it be garbage collected.
		if poolMemoryCounter.Load()+int64(buf.Cap()) < maxLogBufferMemory {
			buf.Reset()
			poolMemoryCounter.Add(int64(buf.Cap()))
			bufPool.Put(buf)
		}
	}()
	jsonHandler := slog.NewJSONHandler(buf, handleOption)
	if len(lo.attrs) > 0 {
		jsonHandler = jsonHandler.WithAttrs(lo.attrs).(*slog.JSONHandler)
	}
	if len(lo.group) > 0 {
		for i := len(lo.group) - 1; i >= 0; i-- {
			jsonHandler = jsonHandler.WithGroup(lo.group[i]).(*slog.JSONHandler)
		}
	}
	if err := jsonHandler.Handle(ctx, rec); err != nil {
		return err
	}
	if !canPretty {
		lo.mu.Lock()
		defer lo.mu.Unlock()
		_, err := io.Copy(os.Stderr, buf)
		return err
	}
	b := pretty.Pretty(buf.Bytes())
	b = pretty.Color(b, nil)
	lo.mu.Lock()
	defer lo.mu.Unlock()
	_, err := io.Copy(os.Stderr, bytes.NewReader(b))
	return err
}

func (lo *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		attrs: append(lo.attrs, attrs...),
		group: lo.group,
		level: lo.level,
	}
}

func (lo *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		attrs: lo.attrs,
		group: append(lo.group, name),
		level: lo.level,
	}
}

func transformProtoToLog(msg proto.Message) slog.Value {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return slog.StringValue("!ERROR:" + err.Error())
	}
	return slog.AnyValue(json.RawMessage(b))
}
