package common

import (
	"testing"
)

func testRegistry(t *testing.T, r Registry) {
	t.Helper()
	assert := Assert(t, true)

	x := int64(1)
	xObj := r.Create(x)

	y, ok := r.Read(xObj).(int64)
	assert(ok)
	assert(x == y)

	y = int64(2)
	r.Update(xObj, y)

	x, ok = r.Read(xObj).(int64)
	assert(ok)
	assert(x == y)

	r.Delete(xObj)
}

func TestRegistryArray(t *testing.T) {
	testRegistry(t, NewRegistryArray())
}

func TestRegistryMap(t *testing.T) {
	testRegistry(t, NewRegistryMap())
}
