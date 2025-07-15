package types

import (
	"database/sql/driver"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// DurationMilli represents a duration that stores milliseconds in the database
// but provides a time.Duration interface for application use.
type DurationMilli struct {
	time.Duration
}

// NewDurationMilli creates a new DurationMilli from a time.Duration value
func NewDurationMilli(d time.Duration) DurationMilli {
	return DurationMilli{Duration: d}
}

// NewDurationMilliFromMillis creates a new DurationMilli from milliseconds
func NewDurationMilliFromMillis(millis int64) DurationMilli {
	return DurationMilli{Duration: time.Duration(millis) * time.Millisecond}
}

// NewDurationMilliFromProto creates a new DurationMilli from a protobuf duration
func NewDurationMilliFromProto(d *durationpb.Duration) DurationMilli {
	if d == nil {
		return DurationMilli{Duration: 0}
	}
	return DurationMilli{Duration: d.AsDuration()}
}

// Scan implements the sql.Scanner interface for reading from the database
func (d *DurationMilli) Scan(value any) error {
	if value == nil {
		d.Duration = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		d.Duration = time.Duration(v) * time.Millisecond
		return nil
	case int:
		d.Duration = time.Duration(v) * time.Millisecond
		return nil
	case []byte:
		// Handle string representation of integer
		var millis int64
		if _, err := fmt.Sscanf(string(v), "%d", &millis); err != nil {
			return fmt.Errorf("cannot scan %T into DurationMilli: %v", value, err)
		}
		d.Duration = time.Duration(millis) * time.Millisecond
		return nil
	case string:
		// Handle string representation of integer
		var millis int64
		if _, err := fmt.Sscanf(v, "%d", &millis); err != nil {
			return fmt.Errorf("cannot scan %T into DurationMilli: %v", value, err)
		}
		d.Duration = time.Duration(millis) * time.Millisecond
		return nil
	default:
		return fmt.Errorf("cannot scan %T into DurationMilli", value)
	}
}

// Value implements the driver.Valuer interface for writing to the database
func (d DurationMilli) Value() (driver.Value, error) {
	return d.Duration.Milliseconds(), nil
}

// Milliseconds returns the duration as milliseconds
func (d DurationMilli) Milliseconds() int64 {
	return d.Duration.Milliseconds()
}

// ToProto converts DurationMilli to protobuf duration
func (d DurationMilli) ToProto() *durationpb.Duration {
	return durationpb.New(d.Duration)
}

// String returns a string representation of the duration
func (d DurationMilli) String() string {
	return d.Duration.String()
}