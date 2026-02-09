package capture

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/lfcypo/scrkun/logger"
)

var imageLogger = logger.New("Image")

func SaveImage(img *image.RGBA, path string) {
	err := os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		imageLogger.Errorf("Failed to create directory: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		imageLogger.Errorf("Failed to save image file to %s: %v", path, err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	jpegOptions := &jpeg.Options{Quality: 50}

	EncodeJPEGFile(file, img, jpegOptions)
}
