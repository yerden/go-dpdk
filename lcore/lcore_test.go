package lcore

import (
	"bytes"
	"testing"
)

func TestReadCpuHexMap(t *testing.T) {
	s := "4f"
	b := bytes.NewBufferString(s)

	cores, err := readCpuHexMap(b)
	a := []int{0, 1, 2, 3, 6}

	if err != nil {
		t.FailNow()
	}
	for i := range a {
		if a[i] != cores[i] {
			t.FailNow()
		}
	}
	if len(a) != len(cores) {
		t.FailNow()
	}
}

func TestReadCpuHexMapErr(t *testing.T) {
	s := "4z"
	b := bytes.NewBufferString(s)

	_, err := readCpuHexMap(b)

	if err != errInvalidMap {
		t.FailNow()
	}
}
