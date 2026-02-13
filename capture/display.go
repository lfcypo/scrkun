package capture

import (
	"sync"

	"github.com/kbinani/screenshot"
	"github.com/lfcypo/scrkun/logger"
)

var displayLogger = logger.New("Display")

var count int

var doOnceGetDisplayCount = &sync.Once{}

func GetDisplayCount() int {
	doOnceGetDisplayCount.Do(getDisplayCount)

	return count
}

func getDisplayCount() {
	count = screenshot.NumActiveDisplays()
	displayLogger.Infof("共有 %d 块显示器", count)
}

func HasDisplay() bool {
	return GetDisplayCount() > 0
}
