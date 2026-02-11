package capture

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"

	"github.com/lfcypo/scrkun/logger"
)

var encodeLogger = logger.New("Capture Encoder")

func EncodeJPEGFile(buffer io.Writer, img *image.RGBA, options *jpeg.Options) {
	err := jpeg.Encode(buffer, img, options)
	if err != nil {
		encodeLogger.Errorf("JPEG 编码失败: %v", err)
	}
}

func RGBAtoBase64JPEG(img *image.RGBA, quality int) string {
	buf := new(bytes.Buffer)

	jpegOptions := &jpeg.Options{Quality: quality}
	EncodeJPEGFile(buf, img, jpegOptions)

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	// data:image/jpeg;base64,{base64_image}
	return "data:image/jpeg;base64," + base64Image
}
