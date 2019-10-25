package flow

import "testing"

func assert(t testing.TB, expected bool, args ...interface{}) {
	if !expected {
		t.Helper()
		t.Fatal(args...)
	}
}

func TestCPattern(t *testing.T) {
	pattern := []Item{
		{Spec: &ItemIPv4{}, Mask: &ItemIPv4{}},
	}

	pat := cPattern(pattern)
	assert(t, pat != nil)

}
