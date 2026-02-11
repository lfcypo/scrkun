package capture

import (
	"errors"
	"image"

	"github.com/kbinani/screenshot"
	"github.com/lfcypo/scrkun/logger"
)

var leastImage *image.RGBA

var captureLogger = logger.New("Capture")

var Skip = errors.New("skip")

func Capture(savePath string) (*image.RGBA, error) {
	if !HasDisplay() {
		captureLogger.Error("找不到显示器")
		err := errors.New("no display found")
		return nil, err
	}

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		captureLogger.Errorf("屏幕截图失败: %v", err)
		return nil, err
	}

	if leastImage == nil {
		leastImage = img
	} else {
		if CompareSimilar(img, leastImage, 90) {
			captureLogger.Debug("相似的屏幕 跳过本次检测")
			return nil, Skip
		}
	}
	captureLogger.Debugf("截取到桌面: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	leastImage = img
	SaveImage(img, savePath)
	return img, nil
}
