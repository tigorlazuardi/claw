package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SlogOption is a function that configures a slog handler.
type SlogOption func(*slogConfig)

// slogConfig holds the configuration for the slog handler.
type slogConfig struct {
	loggerProvider log.LoggerProvider
	name           string
	otelslogOpts   []otelslog.Option
}

// WithLoggerProvider sets the OpenTelemetry logger provider.
// This is required for the handler to work.
func WithLoggerProvider(provider log.LoggerProvider) SlogOption {
	return func(c *slogConfig) {
		c.loggerProvider = provider
	}
}

// WithName sets the logger name for the handler.
// The logger name is used to identify the source of log records.
func WithName(name string) SlogOption {
	return func(c *slogConfig) {
		c.name = name
	}
}

// WithOtelSlogOptions sets additional options for the underlying otelslog.Handler.
// This allows for fine-grained control over the OpenTelemetry slog bridge.
func WithOtelSlogOptions(opts ...otelslog.Option) SlogOption {
	return func(c *slogConfig) {
		c.otelslogOpts = append(c.otelslogOpts, opts...)
	}
}

// SlogHandler creates a slog.Handler that bridges to OpenTelemetry logging
// with special handling for protobuf messages.
//
// When a log field contains a protobuf message, it will be transformed into
// a slog.GroupValue with all the message fields as structured attributes.
//
// Special protobuf well-known types are handled with native slog types:
//   - timestamppb.Timestamp → slog.TimeValue
//   - durationpb.Duration → slog.DurationValue
//
// Example usage:
//
//	provider, err := NewLoggingProvider()
//	if err != nil {
//		// handle error
//	}
//
//	handler := SlogHandler(WithLoggerProvider(provider.LoggerProvider))
//	logger := slog.New(handler)
//
//	// Logging with protobuf message
//	msg := &pb.MyMessage{
//		Field1: "value1", 
//		Field2: 42,
//		Timestamp: timestamppb.Now(),
//		Duration: durationpb.New(time.Second * 30),
//	}
//	logger.Info("Processing message", "proto", msg)
//
// Advanced usage with custom options:
//
//	handler := SlogHandler(
//		WithLoggerProvider(provider.LoggerProvider),
//		WithName("my-service"),
//		WithOtelSlogOptions(otelslog.WithVersion("1.0.0")),
//	)
func SlogHandler(opts ...SlogOption) slog.Handler {
	return SlogHandlerWithName("github.com/tigorlazuardi/claw", opts...)
}

// SlogHandlerWithName creates a slog.Handler with a specific logger name.
//
// The logger name is used to identify the source of log records and can be
// used for filtering and routing in observability systems.
//
// Example usage:
//
//	handler := SlogHandlerWithName("myapp.service", WithLoggerProvider(provider.LoggerProvider))
//	logger := slog.New(handler)
//
// Advanced usage:
//
//	handler := SlogHandlerWithName("myapp.service",
//		WithLoggerProvider(provider.LoggerProvider),
//		WithOtelSlogOptions(otelslog.WithVersion("1.0.0")),
//	)
func SlogHandlerWithName(name string, opts ...SlogOption) slog.Handler {
	config := &slogConfig{
		name: name,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.loggerProvider == nil {
		config.loggerProvider = global.GetLoggerProvider()
	}

	// Prepare otelslog options
	otelslogOpts := []otelslog.Option{
		otelslog.WithLoggerProvider(config.loggerProvider),
	}
	otelslogOpts = append(otelslogOpts, config.otelslogOpts...)

	return &slogHandler{
		handler: otelslog.NewHandler(config.name, otelslogOpts...),
	}
}

// slogHandler wraps the otelslog.Handler to provide protobuf message transformation.
type slogHandler struct {
	handler slog.Handler
}

// Enabled reports whether the handler handles records at the given level.
func (h *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle handles the Record by transforming protobuf messages and forwarding to otelslog.
func (h *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Transform protobuf messages in the record
	transformedRecord := h.transformRecord(record)
	return h.handler.Handle(ctx, transformedRecord)
}

// WithAttrs returns a new Handler whose attributes consist of both the receiver's
// attributes and the arguments.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Transform protobuf messages in attributes
	transformedAttrs := h.transformAttrs(attrs)
	return &slogHandler{
		handler: h.handler.WithAttrs(transformedAttrs),
	}
}

// WithGroup returns a new Handler with the given group appended to the receiver's
// existing groups.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	return &slogHandler{
		handler: h.handler.WithGroup(name),
	}
}

// transformRecord transforms protobuf messages in a slog.Record.
func (h *slogHandler) transformRecord(record slog.Record) slog.Record {
	// Create a new record with transformed attributes
	newRecord := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	// Transform each attribute
	record.Attrs(func(attr slog.Attr) bool {
		transformedAttr := h.transformAttr(attr)
		newRecord.AddAttrs(transformedAttr)
		return true
	})

	return newRecord
}

// transformAttrs transforms protobuf messages in a slice of slog.Attr.
func (h *slogHandler) transformAttrs(attrs []slog.Attr) []slog.Attr {
	transformed := make([]slog.Attr, len(attrs))
	for i, attr := range attrs {
		transformed[i] = h.transformAttr(attr)
	}
	return transformed
}

// transformAttr transforms protobuf messages in a single slog.Attr.
func (h *slogHandler) transformAttr(attr slog.Attr) slog.Attr {
	switch attr.Value.Kind() {
	case slog.KindAny:
		if protoMsg, ok := attr.Value.Any().(proto.Message); ok {
			// Transform protobuf message to group value
			groupValue := h.protoToGroupValue(protoMsg)
			return slog.Attr{
				Key:   attr.Key,
				Value: groupValue,
			}
		}
	case slog.KindGroup:
		// Recursively transform group attributes
		groupAttrs := attr.Value.Group()
		transformedGroupAttrs := h.transformAttrs(groupAttrs)
		return slog.Attr{
			Key:   attr.Key,
			Value: slog.GroupValue(transformedGroupAttrs...),
		}
	}

	return attr
}

// protoToGroupValue converts a protobuf message to a slog.Value containing a group.
func (h *slogHandler) protoToGroupValue(msg proto.Message) slog.Value {
	if msg == nil {
		return slog.GroupValue()
	}

	reflectMsg := msg.ProtoReflect()
	fields := reflectMsg.Descriptor().Fields()

	attrs := make([]slog.Attr, 0, fields.Len())

	// Iterate through all fields in the message
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if !reflectMsg.Has(field) {
			continue // Skip unset fields
		}

		key := string(field.Name())
		value := reflectMsg.Get(field)

		slogValue := h.protoValueToSlogValue(value, field)
		attrs = append(attrs, slog.Attr{
			Key:   key,
			Value: slogValue,
		})
	}

	return slog.GroupValue(attrs...)
}

// protoValueToSlogValue converts a protobuf field value to a slog.Value.
func (h *slogHandler) protoValueToSlogValue(value protoreflect.Value, field protoreflect.FieldDescriptor) slog.Value {
	if field.IsList() {
		return h.protoListToSlogValue(value.List(), field)
	}

	if field.IsMap() {
		return h.protoMapToSlogValue(value.Map(), field)
	}

	return h.protoScalarToSlogValue(value, field)
}

// protoListToSlogValue converts a protobuf list to a slog.Value.
func (h *slogHandler) protoListToSlogValue(list protoreflect.List, field protoreflect.FieldDescriptor) slog.Value {
	values := make([]any, list.Len())
	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		values[i] = h.protoScalarToSlogValue(item, field).Any()
	}
	return slog.AnyValue(values)
}

// protoMapToSlogValue converts a protobuf map to a slog.Value.
func (h *slogHandler) protoMapToSlogValue(protoMap protoreflect.Map, field protoreflect.FieldDescriptor) slog.Value {
	attrs := make([]slog.Attr, 0, protoMap.Len())

	protoMap.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
		keyStr := key.String()
		slogValue := h.protoScalarToSlogValue(value, field.MapValue())
		attrs = append(attrs, slog.Attr{
			Key:   keyStr,
			Value: slogValue,
		})
		return true
	})

	return slog.GroupValue(attrs...)
}

// protoScalarToSlogValue converts a protobuf scalar value to a slog.Value.
// Special well-known types are handled with native slog types:
//   - timestamppb.Timestamp → slog.TimeValue
//   - durationpb.Duration → slog.DurationValue
func (h *slogHandler) protoScalarToSlogValue(value protoreflect.Value, field protoreflect.FieldDescriptor) slog.Value {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return slog.BoolValue(value.Bool())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return slog.IntValue(int(value.Int()))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return slog.Int64Value(value.Int())
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return slog.IntValue(int(value.Uint()))
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return slog.Uint64Value(value.Uint())
	case protoreflect.FloatKind:
		return slog.Float64Value(float64(value.Float()))
	case protoreflect.DoubleKind:
		return slog.Float64Value(value.Float())
	case protoreflect.StringKind:
		return slog.StringValue(value.String())
	case protoreflect.BytesKind:
		return slog.StringValue(string(value.Bytes()))
	case protoreflect.EnumKind:
		enumDesc := field.Enum()
		enumValueDesc := enumDesc.Values().ByNumber(value.Enum())
		if enumValueDesc != nil {
			return slog.StringValue(string(enumValueDesc.Name()))
		}
		return slog.StringValue("unknown")
	case protoreflect.MessageKind:
		// Handle special well-known types first
		if value.Message().IsValid() {
			nestedMsg := value.Message().Interface()
			
			// Handle timestamppb.Timestamp
			if ts, ok := nestedMsg.(*timestamppb.Timestamp); ok {
				return slog.TimeValue(ts.AsTime())
			}
			
			// Handle durationpb.Duration
			if dur, ok := nestedMsg.(*durationpb.Duration); ok {
				return slog.DurationValue(dur.AsDuration())
			}
			
			// For other message types, recursively convert to group
			return h.protoToGroupValue(nestedMsg)
		}
		return slog.StringValue("")
	default:
		// Fallback to string representation
		return slog.StringValue(value.String())
	}
}

// IsProtoMessage checks if a value is a protobuf message.
func IsProtoMessage(v any) bool {
	if v == nil {
		return false
	}

	// Check if it implements proto.Message
	_, ok := v.(proto.Message)
	return ok
}

// ProtoToGroupValue is a utility function to convert a protobuf message to slog.Value containing a group.
// This can be used outside of the slog handler context.
func ProtoToGroupValue(msg proto.Message) slog.Value {
	handler := &slogHandler{}
	return handler.protoToGroupValue(msg)
}
