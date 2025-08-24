package claw

import (
	"cmp"

	"github.com/go-jet/jet/v2/sqlite"
)

// Clamp restricts a value to be within the specified minimum and maximum bounds.
func Clamp[T cmp.Ordered](value, minimum, maximum T) T {
	return max(minimum, min(value, maximum))
}

func toOrderByClause(field sqlite.Expression, desc bool) sqlite.OrderByClause {
	if desc {
		return field.DESC()
	}
	return field.ASC()
}
