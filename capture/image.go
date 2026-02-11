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
		imageLogger.Errorf("创建临时目录失败: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		imageLogger.Errorf("保存文件到 %s 失败: %v", path, err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	jpegOptions := &jpeg.Options{Quality: 50}

	EncodeJPEGFile(file, img, jpegOptions)
}
