package qwalk

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type ItemInfo struct {
	Info os.FileInfo
	Path string // Absolute Path
}

type FilterHandler func(info ItemInfo) bool

//goland:noinspection ALL
func WalkSlice(
	targetDirAbsPaths []string,
	filterHandler FilterHandler,
	workerCount int,
) (results []ItemInfo) {
	resultsMutex := sync.Mutex{}

	Walk(
		targetDirAbsPaths,
		func(info ItemInfo) bool {
			if filterHandler != nil && !filterHandler(info) {
				return false
			}

			resultsMutex.Lock()
			results = append(results, info)
			resultsMutex.Unlock()

			return true
		},
		workerCount,
	)

	return
}

func Walk(
	targetDirAbsPaths []string,
	filterHandler FilterHandler,
	workerCount int,
) {
	var incompleteRequestCount int64

	workRequests := make(chan string)
	bufferRequests := make(chan string)
	done := make(chan struct{})

	for i := 0; i < workerCount; i++ {
		go WalkWorker(
			workRequests,
			bufferRequests,
			filterHandler,
			&incompleteRequestCount,
			done,
		)
	}

	var buffer []string

	for _, targetDirAbsPath := range targetDirAbsPaths {
		if !filepath.IsAbs(targetDirAbsPath) {
			targetDirAbsPath, _ = filepath.Abs(targetDirAbsPath)
		}

		buffer = append(buffer, targetDirAbsPath)

		atomic.AddInt64(&incompleteRequestCount, 1)
	}

	for {
		if len(buffer) > 0 {
			select {
			case workRequests <- buffer[0]:
				buffer = buffer[1:]
			case bufferRequest := <-bufferRequests:
				buffer = append(buffer, bufferRequest)
			}
		} else {
			select {
			case bufferRequest := <-bufferRequests:
				buffer = append(buffer, bufferRequest)
			case <-done:
				goto exitFor
			}
		}
	}
exitFor:
	close(workRequests)
	close(bufferRequests)
	close(done)
}

func WalkWorker(
	workRequests chan string,
	bufferRequests chan string,
	filterHandler FilterHandler,
	incompleteRequestCount *int64,
	done chan struct{},
) {
	for {
		request, ok := <-workRequests

		if !ok {
			return
		}

		f, _ := os.Open(request)
		fsItems, err := f.Readdir(1)

		for err != io.EOF && len(fsItems) > 0 {
			fsItem := fsItems[0]
			absPath := filepath.Join(request, fsItem.Name())

			fsi := ItemInfo{
				fsItem,
				absPath,
			}

			allowRequest := true

			if filterHandler != nil {
				_allowRequest := filterHandler(fsi)

				if !_allowRequest {
					fsItems, err = f.Readdir(1)

					continue
				}

				allowRequest = _allowRequest
			}

			if allowRequest && fsItem.IsDir() {
				atomic.AddInt64(incompleteRequestCount, 1)

				bufferRequests <- absPath
			}

			fsItems, err = f.Readdir(1)
		}

		_ = f.Close()

		if atomic.AddInt64(incompleteRequestCount, -1) == 0 {
			done <- struct{}{}
		}
	}
}
