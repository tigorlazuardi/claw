package types

import (
	"database/sql/driver"
	"fmt"
)

// Bool represents a boolean value that stores 0/1 in the database
// but provides a normal bool interface for application use.
type Bool bool

// NewBool creates a new Bool from a bool value
func NewBool(b bool) Bool {
	return Bool(b)
}

// NewBoolFromPointer creates a new Bool from a *bool value
// Returns false if the pointer is nil
func NewBoolFromPointer(b *bool) Bool {
	if b == nil {
		return Bool(false)
	}
	return Bool(*b)
}

// Scan implements the sql.Scanner interface for reading from the database
func (b *Bool) Scan(value any) error {
	if value == nil {
		*b = Bool(false)
		return nil
	}

	switch v := value.(type) {
	case int64:
		*b = Bool(v != 0)
		return nil
	case int:
		*b = Bool(v != 0)
		return nil
	case bool:
		*b = Bool(v)
		return nil
	case []byte:
		// Handle string representation of integer
		var i int64
		if _, err := fmt.Sscanf(string(v), "%d", &i); err != nil {
			return fmt.Errorf("cannot scan %T into Bool: %v", value, err)
		}
		*b = Bool(i != 0)
		return nil
	case string:
		// Handle string representation of integer
		var i int64
		if _, err := fmt.Sscanf(v, "%d", &i); err != nil {
			return fmt.Errorf("cannot scan %T into Bool: %v", value, err)
		}
		*b = Bool(i != 0)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Bool", value)
	}
}

// Value implements the driver.Valuer interface for writing to the database
func (b Bool) Value() (driver.Value, error) {
	if b {
		return int64(1), nil
	}
	return int64(0), nil
}

// Bool returns the bool value
func (b Bool) Bool() bool {
	return bool(b)
}

// Pointer returns a pointer to the bool value
func (b Bool) Pointer() *bool {
	v := bool(b)
	return &v
}

// String returns a string representation of the boolean
func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}