package core

import (
	"strconv"
	"strings"
)

// IndexOutOfRange reports whether a selector's index names no element among
// count matches. `index: 1` with a single match used to fall back to the first
// match, so "there is no second row" failed and a tap hit the wrong row (#188).
// Negative indexes count from the end. An empty or non-numeric index is not
// out of range: those keep their existing handling.
func IndexOutOfRange(count int, index string) bool {
	i, err := strconv.Atoi(strings.TrimSpace(index))
	if err != nil {
		return false
	}
	if i < 0 {
		i += count
	}
	return i < 0 || i >= count
}
