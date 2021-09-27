package experimental_test

import (
	"github.com/nikosgram/qwalk/experimental"
	"runtime"
	"testing"
)

func BenchmarkWalk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := experimental.Walk(
			[]string{"."},
			int64(runtime.NumCPU()),
			func(chunk []experimental.ItemInfo) {},
		)

		if err != nil {
			b.Error(err)
		}
	}
}
