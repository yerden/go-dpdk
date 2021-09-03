package util

import (
	"fmt"
	"sort"
	"strings"
)

// LcoresList represents a list of unsigned ints, e.g. lcores.
// It implements sort.Interface.
type LcoresList []uint

func (list LcoresList) Len() int {
	return len(list)
}

func (list LcoresList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list LcoresList) Less(i, j int) bool {
	return list[i] < list[j]
}

var _ sort.Interface = LcoresList{}

// Sort sorts list of unsigned ints.
func (list LcoresList) Sort() {
	sort.Sort(list)
}

// Dup allocates new list and copies list's contents into it.
func (list LcoresList) Dup() LcoresList {
	newList := append(LcoresList{}, list...)
	newList.Sort()
	return newList
}

func (list LcoresList) String() string {
	if list.Len() == 0 {
		return ""
	}

	list = list.Dup()
	prev := 0
	next := 0
	var s []string

	for i := 1; i < list.Len(); i++ {
		if list[i] == list[i-1] {
			continue
		}

		if list[i]-list[i-1] == 1 {
			next = i
			continue
		}

		s = append(s, list.rangeStr(prev, next))
		prev = i
		next = i
	}

	s = append(s, list.rangeStr(prev, next))
	return strings.Join(s, ",")
}

func (list LcoresList) rangeStr(prev, next int) string {
	if prev == next {
		return fmt.Sprintf("%d", list[prev])
	}
	return fmt.Sprintf("%d-%d", list[prev], list[next])
}

// Equal returns true if list and other are identical.
func (list LcoresList) Equal(other LcoresList) bool {
	list = list.Dup()
	other = other.Dup()

	if list.Len() != other.Len() {
		return false
	}

	for i := range list {
		if list[i] != other[i] {
			return false
		}
	}

	return true
}
