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
			nil,
			runtime.NumCPU(),
		)
	}
}

func BenchmarkWalkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = qwalk.WalkSlice(
			[]string{"."},
			nil,
			runtime.NumCPU(),
		)
	}
}
