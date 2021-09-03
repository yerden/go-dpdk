package util

import "testing"

func assert(t testing.TB, expected bool, args ...interface{}) {
	if !expected {
		t.Helper()
		t.Fatal(args...)
	}
}

func TestLcoresBasic(t *testing.T) {
	list := LcoresList{2, 4, 5, 6, 7, 8, 1}
	sorted := LcoresList{1, 2, 4, 5, 6, 7, 8}

	assert(t, sorted.Equal(list))
	assert(t, list.Equal(sorted))

	list.Sort()
}

var testCases = []struct {
	LcoresList
	String string
}{
	{LcoresList{2, 4, 5, 6, 7, 8, 1}, "1-2,4-8"},
	{LcoresList{2, 4, 5, 6}, "2,4-6"},
	{LcoresList{4, 5, 2, 4}, "2,4-5"},
	{LcoresList{4, 5, 6, 7, 7}, "4-7"},
	{LcoresList{4, 5, 7, 8, 7}, "4-5,7-8"},
	{LcoresList{4, 5, 6, 7, 9}, "4-7,9"},
	{LcoresList{4, 6, 7, 8, 19}, "4,6-8,19"},
}

func TestLcoresString(t *testing.T) {
	for i := range testCases {
		list := testCases[i].LcoresList
		s := testCases[i].String
		assert(t, list.String() == s, list.String(), s)
	}
}
