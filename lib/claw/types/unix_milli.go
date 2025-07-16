package types

import (
	"database/sql/driver"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// UnixMilli represents a timestamp that stores Unix milliseconds in the database
// but provides a time.Time interface for application use.
type UnixMilli struct {
	time.Time
}

// NewUnixMilli creates a new UnixMilli from a time.Time value
func NewUnixMilli(t time.Time) UnixMilli {
	return UnixMilli{Time: t}
}

// NewUnixMilliFromUnix creates a new UnixMilli from Unix milliseconds
func NewUnixMilliFromUnix(millis int64) UnixMilli {
	return UnixMilli{Time: time.UnixMilli(millis)}
}

// NewUnixMilliFromProto creates a new UnixMilli from a protobuf timestamp
func NewUnixMilliFromProto(ts *timestamppb.Timestamp) UnixMilli {
	if ts == nil {
		return UnixMilli{Time: time.Time{}}
	}
	return UnixMilli{Time: ts.AsTime()}
}

// Now returns the current time as UnixMilli
func Now() UnixMilli {
	return UnixMilli{Time: time.Now()}
}

// UnixMilliNow returns the current time as UnixMilli (alias for Now)
func UnixMilliNow() UnixMilli {
	return Now()
}

// Scan implements the sql.Scanner interface for reading from the database
func (u *UnixMilli) Scan(value any) error {
	if value == nil {
		u.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case int64:
		u.Time = time.UnixMilli(v)
		return nil
	case int:
		u.Time = time.UnixMilli(int64(v))
		return nil
	case []byte:
		// Handle string representation of integer
		var millis int64
		if _, err := fmt.Sscanf(string(v), "%d", &millis); err != nil {
			return fmt.Errorf("cannot scan %T into UnixMilli: %v", value, err)
		}
		u.Time = time.UnixMilli(millis)
		return nil
	case string:
		// Handle string representation of integer
		var millis int64
		if _, err := fmt.Sscanf(v, "%d", &millis); err != nil {
			return fmt.Errorf("cannot scan %T into UnixMilli: %v", value, err)
		}
		u.Time = time.UnixMilli(millis)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into UnixMilli", value)
	}
}

// Value implements the driver.Valuer interface for writing to the database
func (u UnixMilli) Value() (driver.Value, error) {
	if u.Time.IsZero() {
		return nil, nil
	}
	return u.Time.UnixMilli(), nil
}

// UnixMilli returns the Unix milliseconds representation
func (u UnixMilli) UnixMilli() int64 {
	return u.Time.UnixMilli()
}

// ToProto converts UnixMilli to protobuf timestamp
func (u UnixMilli) ToProto() *timestamppb.Timestamp {
	if u.Time.IsZero() {
		return nil
	}
	return timestamppb.New(u.Time)
}

// String returns a string representation of the timestamp
func (u UnixMilli) String() string {
	if u.Time.IsZero() {
		return "0"
	}
	return u.Time.Format(time.RFC3339)
}
