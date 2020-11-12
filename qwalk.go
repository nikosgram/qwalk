package main

import (
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
)

type TFsItemInfo struct {
	Info    os.FileInfo
	AbsPath string
	RelPath string
}

func GenerateFsItemsSlice(
	targetDirAbsPath string,
	workerCount int,
	filterHandler func(fsItemInfo TFsItemInfo) (bool, bool),
) []TFsItemInfo {
	fsItemChan := make(chan TFsItemInfo)

	go func() {
		GenerateFsItems(targetDirAbsPath, fsItemChan, workerCount, filterHandler)

		close(fsItemChan)
	}()

	returnedFsItems := make([]TFsItemInfo, 0, (1<<16)-1)

	for {
		fsItemInfo, ok := <-fsItemChan

		if !ok {
			return returnedFsItems
		}

		returnedFsItems = append(returnedFsItems, fsItemInfo)
	}
}

func GenerateFsItems(
	targetDirAbsPath string,
	results chan TFsItemInfo,
	workerCount int,
	filterHandler func(fsItemInfo TFsItemInfo) (bool, bool),
) {
	var incompleteRequestCount int64

	workRequests := make(chan string)
	bufferRequests := make(chan string)
	done := make(chan struct{})

	for i := 0; i < workerCount; i++ {
		go DirListingWorker(
			targetDirAbsPath,
			workRequests,
			bufferRequests,
			results,
			&incompleteRequestCount,
			done,
			filterHandler,
		)
	}

	var buffer []string

	buffer = append(buffer, targetDirAbsPath)

	atomic.AddInt64(&incompleteRequestCount, 1)

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
}

func DirListingWorker(
	targetDirAbsPath string,
	workRequests chan string,
	bufferRequests chan string,
	results chan TFsItemInfo,
	incompleteRequestCount *int64,
	done chan struct{},
	filterHandler func(fsItemInfo TFsItemInfo) (bool, bool),
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
			relPath, _ := filepath.Rel(targetDirAbsPath, absPath)

			fsi := TFsItemInfo{
				fsItem,
				absPath,
				relPath,
			}

			allowResult := true

			if filterHandler != nil {
				allowRequest, _allowResult := filterHandler(fsi)

				if !allowRequest {
					fsItems, err = f.Readdir(1)

					continue
				}

				allowResult = _allowResult
			}

			if allowResult {
				results <- fsi
			}

			if fsItem.IsDir() {
				atomic.AddInt64(incompleteRequestCount, 1)

				bufferRequests <- absPath
			}

			fsItems, err = f.Readdir(1)
		}

		_ = f.Close()

		atomic.AddInt64(incompleteRequestCount, -1)

		if atomic.LoadInt64(incompleteRequestCount) == 0 {
			done <- struct{}{}
		}
	}
}
