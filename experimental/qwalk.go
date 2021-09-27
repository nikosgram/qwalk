package experimental

import (
	"os"
)

type ItemInfo struct {
	Info *os.FileInfo
	Path string // Absolute Path
}

type ResultHandler func(chunk []ItemInfo)

func Walk(
	sourceDirs []string,
	threadsCount int64,
	handler ResultHandler,
) (
	err error,
) {
	return
}
