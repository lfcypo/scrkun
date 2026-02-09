package capture

import (
	"github.com/kbinani/screenshot"
	"github.com/lfcypo/scrkun/logger"
)

var displayLogger = logger.New("Display")

func GetDisplayCount() int {
	count := screenshot.NumActiveDisplays()
	displayLogger.Infof("Displays count: %d", count)
	return count
}

func HasDisplay() bool {
	return GetDisplayCount() > 0
}
