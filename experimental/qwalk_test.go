package experimental_test

import (
	"github.com/nikosgram/qwalk/experimental"
	"runtime"
	"testing"
)

func TestWalk(t *testing.T) {
	err := experimental.Walk(
		[]string{"."},
		int64(runtime.NumCPU()),
		func(chunk []experimental.ItemInfo) {
			for _, info := range chunk {
				t.Log(info.Path)
			}
		},
	)

	if err != nil {
		t.Error(err)
	}
}
