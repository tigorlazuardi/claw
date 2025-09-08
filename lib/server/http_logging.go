package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture response information
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

// Write captures the number of bytes written
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// HTTPLoggingMiddleware creates a logging middleware for HTTP requests.
// It logs Method, Path, Response Bytes (human friendly), Status Code, and Duration
// in a single message field using sprintf formatting.
func HTTPLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap the response writer to capture information
			rw := &responseWriter{ResponseWriter: w}
			
			// Process the request
			next.ServeHTTP(rw, r)
			
			duration := time.Since(start)
			
			// Ensure status code is set
			if rw.statusCode == 0 {
				rw.statusCode = http.StatusOK
			}
			
			// Format response bytes in human-friendly format
			responseSize := formatBytes(rw.bytes)
			
			// Log in a single message using sprintf
			message := fmt.Sprintf("%s %s - %d - %s - %v",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				responseSize,
				duration,
			)
			
			logger.Info(message)
		})
	}
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}