package core

import "testing"

// An index past the matches names no element (#188): `index: 1` with one match
// is out of range, where it used to fall back to the first match.
func TestIndexOutOfRange(t *testing.T) {
	for _, tc := range []struct {
		count int
		index string
		want  bool
	}{
		{1, "0", false},
		{1, "1", true}, // "there is no second row"
		{3, "2", false},
		{3, "3", true},
		{3, "-1", false}, // last
		{3, "-3", false}, // first
		{3, "-4", true},
		{0, "0", true},
		{2, "", false},    // no index: unchanged handling
		{2, "abc", false}, // not a number: unchanged handling
		{2, " 1 ", false},
	} {
		if got := IndexOutOfRange(tc.count, tc.index); got != tc.want {
			t.Errorf("IndexOutOfRange(%d, %q) = %v, want %v", tc.count, tc.index, got, tc.want)
		}
	}
}
