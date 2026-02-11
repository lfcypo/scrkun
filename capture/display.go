package capture

import (
	"github.com/kbinani/screenshot"
	"github.com/lfcypo/scrkun/logger"
)

var displayLogger = logger.New("Display")

func GetDisplayCount() int {
	count := screenshot.NumActiveDisplays()
	displayLogger.Infof("共有 %d 块显示器", count)
	return count
}

func HasDisplay() bool {
	return GetDisplayCount() > 0
}
