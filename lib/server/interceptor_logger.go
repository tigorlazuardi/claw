package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// LoggingInterceptor creates a ConnectRPC interceptor for comprehensive logging.
//
// For unary RPCs, it logs:
// - Request procedure, headers, and body
// - Response headers and body
// - Round trip duration
// - Success/failure status
//
// For streaming RPCs, it logs:
// - Connection opened event with headers
// - Connection closed event with duration
// - Does not log message bodies to avoid performance impact
//
// The interceptor uses the provided slog.Logger instance and includes
// structured logging with context information.
func LoggingInterceptor(logger *slog.Logger) connect.Interceptor {
	return &loggerInterceptor{logger: logger}
}

type loggerInterceptor struct {
	logger *slog.Logger
}

// WrapUnary implements logging for unary RPC calls
func (c *loggerInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		procedure := req.Spec().Procedure

		// Extract headers
		headers := make(logHeader)
		maps.Copy(headers, req.Header())

		// Extract request body
		var reqBody json.RawMessage
		if msg, ok := req.Any().(proto.Message); ok {
			if bodyBytes, err := protojson.Marshal(msg); err == nil {
				reqBody = json.RawMessage(bodyBytes)
			}
		}

		// Call the actual RPC
		resp, err := next(ctx, req)
		duration := time.Since(start)

		// Extract response body if successful
		var respBody any
		var responseHeaders map[string][]string
		if err == nil && resp != nil {
			respBody = resp.Any()
			// Extract response headers
			responseHeaders = make(logHeader)
			maps.Copy(responseHeaders, resp.Header())
		}

		// Log the complete request/response cycle
		if err == nil {
			msg := fmt.Sprintf("RPC %s - ok - %s", procedure, duration)
			c.logger.InfoContext(ctx, msg)
			c.logger.DebugContext(ctx, msg,
				slog.String("type", "unary_rpc"),
				slog.String("procedure", procedure),
				slog.Duration("duration", duration),
				slog.Any("request_headers", headers),
				slog.Any("request_body", reqBody),
				slog.Any("response_headers", responseHeaders),
				slog.Any("response_body", respBody),
				slog.String("status", "success"),
			)
		} else {
			code := connect.CodeInternal
			if e := (&connect.Error{}); errors.As(err, &e) {
				code = e.Code()
			}
			msg := fmt.Sprintf("RPC %s - %s - %s", procedure, code, duration)
			c.logger.ErrorContext(ctx, msg,
				slog.String("type", "unary_rpc"),
				slog.String("procedure", procedure),
				slog.Duration("duration", duration),
				slog.Any("request_headers", headers),
				slog.Any("request_body", reqBody),
				slog.String("status", "error"),
				slog.String("error", err.Error()),
			)
		}

		return resp, err
	}
}

// WrapStreamingClient implements logging for streaming client calls
func (c *loggerInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, spec)
		return &loggingStreamingClientConn{
			StreamingClientConn: conn,
			logger:              c.logger,
			procedure:           spec.Procedure,
			start:               time.Now(),
		}
	}
}

// WrapStreamingHandler implements logging for streaming handler calls
func (c *loggerInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		start := time.Now()
		procedure := conn.Spec().Procedure

		// Extract headers from the streaming connection
		headers := make(logHeader)
		maps.Copy(headers, conn.RequestHeader())

		// Log connection opened
		openMsg := fmt.Sprintf("Streaming RPC %s connection open", procedure)
		c.logger.InfoContext(ctx, openMsg,
			slog.String("type", "streaming_rpc"),
			slog.String("procedure", procedure),
			slog.String("event", "connection_opened"),
			slog.Any("headers", headers),
		)

		// Call the actual streaming RPC
		err := next(ctx, conn)
		duration := time.Since(start)

		// Extract response headers
		responseHeaders := make(logHeader)
		maps.Copy(responseHeaders, conn.ResponseHeader())

		// Log connection closed
		if err == nil {
			closeMsg := fmt.Sprintf("Streaming RPC %s connection closed", procedure)
			c.logger.InfoContext(ctx, closeMsg,
				slog.String("type", "streaming_rpc"),
				slog.String("procedure", procedure),
				slog.String("event", "connection_closed"),
				slog.Duration("duration", duration),
				slog.Any("response_headers", responseHeaders),
				slog.String("status", "success"),
			)
		} else {
			closeMsg := fmt.Sprintf("Streaming RPC %s connection closed with error", procedure)
			c.logger.ErrorContext(ctx, closeMsg,
				slog.String("type", "streaming_rpc"),
				slog.String("procedure", procedure),
				slog.String("event", "connection_closed"),
				slog.Duration("duration", duration),
				slog.Any("response_headers", responseHeaders),
				slog.String("status", "error"),
				slog.String("error", err.Error()),
			)
		}

		return err
	}
}

type logHeader map[string][]string

func (lo logHeader) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, len(lo))
	for k, v := range lo {
		if len(v) == 1 {
			attrs = append(attrs, slog.String(k, v[0]))
		} else if len(v) > 1 {
			attrs = append(attrs, slog.Any(k, v))
		}
	}
	return slog.GroupValue(attrs...)
}

// loggingStreamingClientConn wraps a streaming client connection for logging
type loggingStreamingClientConn struct {
	connect.StreamingClientConn
	logger    *slog.Logger
	procedure string
	start     time.Time
}

// CloseRequest logs when the streaming client connection is closed
func (l *loggingStreamingClientConn) CloseRequest() error {
	err := l.StreamingClientConn.CloseRequest()
	duration := time.Since(l.start)

	// Extract response headers
	responseHeaders := make(logHeader)
	maps.Copy(responseHeaders, l.ResponseHeader())

	if err == nil {
		l.logger.Info("Streaming RPC client connection closed",
			slog.String("type", "streaming_rpc_client"),
			slog.String("procedure", l.procedure),
			slog.String("event", "connection_closed"),
			slog.Duration("duration", duration),
			slog.Any("response_headers", responseHeaders),
			slog.String("status", "success"),
		)
	} else {
		l.logger.Error("Streaming RPC client connection closed with error",
			slog.String("type", "streaming_rpc_client"),
			slog.String("procedure", l.procedure),
			slog.String("event", "connection_closed"),
			slog.Duration("duration", duration),
			slog.Any("response_headers", responseHeaders),
			slog.String("status", "error"),
			slog.String("error", err.Error()),
		)
	}

	return err
}
