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
		captureLogger.Error("No display found")
		err := errors.New("no display found")
		return nil, err
	}

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		captureLogger.Errorf("Failed to capture screen: %v", err)
		return nil, err
	}

	if leastImage == nil {
		leastImage = img
	} else {
		if CompareSimilar(img, leastImage, 90) {
			captureLogger.Debug("Image is similar to the last one, skipping")
			return nil, Skip
		}
	}
	captureLogger.Debugf("Captured image: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	leastImage = img
	SaveImage(img, savePath)
	return img, nil
}
