package claw

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/tigorlazuardi/claw/lib/claw/config"
)

// StallReader wraps an io.Reader and monitors for stalled reads based on the provided StallMonitor configuration.
// It tracks read progress and can detect when read operations are consistently slow or stalled.
type StallReader struct {
	ctx        context.Context
	source     io.Reader
	monitor    config.StallMonitor
	totalBytes int64
	debounce   *time.Timer
	start      time.Time
	stallError atomic.Pointer[StallError]
}

// NewStallReader creates a new StallReader that wraps the provided io.Reader with stall monitoring
// using the given StallMonitor configuration.
func NewStallReader(ctx context.Context, reader io.Reader, monitor config.StallMonitor) *StallReader {
	return &StallReader{
		ctx:     ctx,
		source:  reader,
		monitor: monitor,
		start:   time.Now(),
	}
}

// Read implements io.Reader interface and monitors for stalled reads.
// It tracks the time between reads and the amount of data transferred to detect stalls.
// Returns ErrStallDetected immediately when a stall condition is detected.
func (sr *StallReader) Read(p []byte) (n int, err error) {
	if !sr.monitor.Enabled {
		return sr.source.Read(p)
	}

	// Idempotency checks.
	//
	// If an error has already been recorded, return it immediately without attempting further reads.

	// Check if context is done
	if sr.ctx.Err() != nil {
		return 0, sr.ctx.Err()
	}

	// Check if already stalled
	if err := sr.stallError.Load(); err != nil {
		return 0, err
	}

	type result struct {
		n   int
		err error
	}

	readCh := make(chan result, 1)
	go func() {
		n, err := sr.source.Read(p)
		readCh <- result{n: n, err: err}
	}()
	noDataSentTimer := time.NewTimer(sr.monitor.NoDataReceivedDuration)
	select {
	case <-sr.ctx.Done():
		return 0, sr.ctx.Err()
	case <-noDataSentTimer.C:
		noDataSentTimer.Stop()
		err := &StallError{Cause: fmt.Sprintf("no single bytes received for %s", sr.monitor.NoDataReceivedDuration)}
		sr.stallError.Store(err)
		return 0, err
	case res := <-readCh:
		n, err = res.n, res.err
	}

	sr.totalBytes += int64(n)

	// Calculate current speed
	elapsed := time.Since(sr.start)
	if elapsed <= 0 {
		elapsed = time.Millisecond // Prevent division by zero
	}
	currentSpeed := sr.totalBytes / int64(elapsed/time.Second)
	threshold := sr.monitor.Speed

	// Check if speed is below threshold for the configured duration
	if currentSpeed < int64(threshold) {
		// Start a fresh debounce timer if not already started
		// otherwise let the existing timer run
		// to avoid resetting the countdown on every read
		// that is below the threshold.
		if sr.debounce == nil {
			sr.debounce = time.AfterFunc(sr.monitor.SpeedDuration, func() {
				sr.stallError.Store(&StallError{
					Cause: fmt.Sprintf("speed below %s for %s (current: %s)", threshold, sr.monitor.SpeedDuration, humanize.Bytes(uint64(currentSpeed))),
				})
			})
		}
	} else {
		// If speed is back above threshold, stop the debounce timer
		if sr.debounce != nil {
			sr.debounce.Stop()
			sr.debounce = nil
		}
	}

	return n, err
}

// IsStalled returns true if a stall condition has been detected based on the monitor configuration.
func (sr *StallReader) IsStalled() bool {
	return sr.stallError.Load() != nil
}

// StallError returns the StallError if a stall condition has been detected, or nil otherwise.
func (sr *StallReader) StallError() *StallError {
	return sr.stallError.Load()
}

// TotalBytesRead returns the total number of bytes read through this reader.
func (sr *StallReader) TotalBytesRead() int64 {
	return sr.totalBytes
}

// StallError represents an error condition when a stall is detected.
type StallError struct {
	Cause string
}

func (e StallError) Error() string {
	return fmt.Sprintf("stall error: %s", e.Cause)
}
