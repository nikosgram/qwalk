package qwalk_test

import (
	"github.com/nikosgram/qwalk"
	"runtime"
	"testing"
)

func TestWalk(t *testing.T) {
	qwalk.Walk(
		[]string{"."},
		func(info qwalk.ItemInfo) (bool, bool) {
			return true, true
		},
		runtime.NumCPU(),
		func(info qwalk.ItemInfo) {
			t.Log(info.Path)
		},
	)
}

func TestWalkSlice(t *testing.T) {
	slice := qwalk.WalkSlice(
		[]string{"."},
		func(info qwalk.ItemInfo) (bool, bool) {
			return true, true
		},
		runtime.NumCPU(),
	)

	for _, info := range slice {
		t.Log(info.Path)
	}
}
