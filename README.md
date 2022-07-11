# qwalk [![Go](https://github.com/nikosgram/qwalk/actions/workflows/go.yml/badge.svg)](https://github.com/nikosgram/qwalk/actions/workflows/go.yml) [![CodeQL](https://github.com/nikosgram/qwalk/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/nikosgram/qwalk/actions/workflows/codeql-analysis.yml)

Golang fastest directory walking method

Based on https://github.com/feyrob/godirlist

## Examples

There are a few examples in the qwalk_test.go and qwalk_bench_test.go files if you want to see a very simplified version of what you can do with qwalk.

```go
package main

import (
	"fmt"
	"runtime"
	
	"github.com/nikosgram/qwalk"
)

func main() {
	// print all no-directory items
	qwalk.Walk(
		[]string{"."},
		func(info qwalk.ItemInfo) bool {
			// print only from no-dir items
			if !info.Info.IsDir() {
				fmt.Println(info.Path)
			}

			// allow dir-listing on all directories
			return true
		},
		runtime.NumCPU(),
	)
}
```
