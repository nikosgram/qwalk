package qwalk_test

import (
	"github.com/nikosgram/qwalk"
	"runtime"
	"testing"
)

func TestWalk(t *testing.T) {
	qwalk.Walk(
		[]string{"."},
		func(info qwalk.ItemInfo) bool {
			t.Log(info.Path)

			return true
		},
		runtime.NumCPU(),
	)
}

func TestWalkSlice(t *testing.T) {
	slice := qwalk.WalkSlice(
		[]string{"."},
		nil,
		runtime.NumCPU(),
	)

	for _, info := range slice {
		t.Log(info.Path)
	}
}
