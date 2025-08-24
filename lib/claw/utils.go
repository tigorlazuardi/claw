package claw

import (
	"cmp"

	"github.com/go-jet/jet/v2/sqlite"
)

// Clamp restricts a value to be within the specified minimum and maximum bounds.
func Clamp[T cmp.Ordered](value, minimum, maximum T) T {
	return max(minimum, min(value, maximum))
}

// Defer returns the value pointed to by the given pointer.
//
// If the pointer is nil, it returns the zero value of the type T.
func Defer[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}
	return *value
}

func toOrderByClause(field sqlite.Expression, desc bool) sqlite.OrderByClause {
	if desc {
		return field.DESC()
	}
	return field.ASC()
}
