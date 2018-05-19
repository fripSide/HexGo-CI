package hexgo

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"sync/atomic"
	"time"
	"sync"
	"path/filepath"
	//"strings"
)


/*
持续集成：
1. debug时的LiveReload功能。
https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2
*/

// 检测根目录下全部文件是否被修改
func checkFiles(path string, fileMap map[string]time.Time) bool {
	change := false
	traverseFile(path, func(filePath string, info os.FileInfo) {
		//fmt.Println(filePath)
		if _, ok := fileMap[filePath]; ok {
			modTi := fileMap[filePath]
			nowTi := info.ModTime()
			if nowTi.After(modTi.Add(time.Second)) {
				change = true
				fileMap[filePath] = nowTi
				fmt.Println(filePath, modTi, info.ModTime())
			}
		}
	})
	return change
}

func LiveReload(h *HexGo) {
	var once = &sync.Once{}
	isStartFlag := int32(0)
	isStop := int32(0)
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	startServer := func() {
		atomic.StoreInt32(&isStartFlag, 1)
		once.Do(func() {

			h.Run()
		})
	}
	stopServer := func() {
		<- stopSignal
		h.Stop()
		atomic.CompareAndSwapInt32(&isStop, 0, 1)
	}
	path, _ := os.Getwd()
	path = filepath.Join(path, "/")
	fileMap := make(map[string]time.Time)
	updateCache := func() {
		traverseFile(path, func(filePath string, info os.FileInfo) {
			fileMap[filePath] = info.ModTime()
		})
	}
	updateCache()
	go stopServer()
	for {
		if atomic.LoadInt32(&isStartFlag) == 0 {
			go startServer()
		} else {
			res := checkFiles(path, fileMap)
			if res {
				h.Stop()
				atomic.StoreInt32(&isStartFlag, 0)
				once = &sync.Once{}
				fmt.Println("live reload")
				h.SetupFunc()
				updateCache()
			}
			time.Sleep(time.Millisecond * 100)
		}

		if atomic.LoadInt32(&isStop) > 0 {
			return
		}
	}
}