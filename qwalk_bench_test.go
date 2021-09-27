package qwalk_test

import (
	"github.com/nikosgram/qwalk"
	"runtime"
	"testing"
)

func BenchmarkWalk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qwalk.Walk(
			[]string{"."},
			func(info qwalk.ItemInfo) (bool, bool) {
				return true, true
			},
			runtime.NumCPU(),
			func(info qwalk.ItemInfo) {},
		)
	}
}

func BenchmarkWalkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = qwalk.WalkSlice(
			[]string{"."},
			func(info qwalk.ItemInfo) (bool, bool) {
				return true, true
			},
			runtime.NumCPU(),
		)
	}
}
