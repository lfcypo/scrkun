package capture

import (
	"image"
	"math"
)

type ImageComparator struct {
	HashThreshold float64 // 哈希相似度阈值
	ColorWeight   float64 // 颜色直方图权重
	SSIMWeight    float64 // SSIM 权重
	HashWeight    float64 // 哈希权重
	MinImageSize  int     // 最小图像尺寸
}

func DefaultComparator() *ImageComparator {
	return &ImageComparator{
		HashThreshold: 0.8,
		ColorWeight:   0.3,
		SSIMWeight:    0.4,
		HashWeight:    0.3,
		MinImageSize:  16,
	}
}

var hammingTable [256]int

func init() {
	for i := 0; i < 256; i++ {
		count := 0
		n := i
		for n > 0 {
			count++
			n &= n - 1
		}
		hammingTable[i] = count
	}
}

func CompareSimilar(img1, img2 *image.RGBA, threshold int) bool {
	comparator := DefaultComparator()
	comparator.HashThreshold = float64(threshold) / 100.0

	isSimilar, _ := comparator.Compare(img1, img2)
	return isSimilar
}

func (c *ImageComparator) Compare(img1, img2 *image.RGBA) (bool, float64) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return false, 0
	}

	width := bounds1.Dx()
	height := bounds1.Dy()

	if width < c.MinImageSize || height < c.MinImageSize {
		hash1 := dHash(img1)
		hash2 := dHash(img2)
		similarity := hashSimilarityFast(hash1, hash2)
		return similarity >= c.HashThreshold, similarity
	}

	hashSim := hashSimilarityFast(dHash(img1), dHash(img2))
	colorSim := colorHistogramSimilarity(img1, img2)
	ssimVal := ssim(img1, img2)

	totalSimilarity := hashSim*c.HashWeight +
		colorSim*c.ColorWeight +
		ssimVal*c.SSIMWeight

	totalSimilarity = math.Min(1.0, math.Max(0.0, totalSimilarity))

	return totalSimilarity >= c.HashThreshold, totalSimilarity
}

func aHash(img *image.RGBA) uint64 {
	const size = 8
	var pixels [size][size]float64

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			srcX := (x * width) / size
			srcY := (y * height) / size
			c := img.RGBAAt(srcX+bounds.Min.X, srcY+bounds.Min.Y)

			// 灰度转换
			gray := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
			pixels[y][x] = gray
		}
	}

	var sum float64
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sum += pixels[y][x]
		}
	}
	avg := sum / (size * size)

	var hash uint64
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			hash <<= 1
			if pixels[y][x] > avg {
				hash |= 1
			}
		}
	}

	return hash
}

func dHash(img *image.RGBA) uint64 {
	const width = 9
	const height = 8

	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	grayPixels := make([]float64, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := (x * imgWidth) / width
			srcY := (y * imgHeight) / height

			c := img.RGBAAt(
				srcX+bounds.Min.X,
				srcY+bounds.Min.Y,
			)

			gray := 0.299*float64(c.R) +
				0.587*float64(c.G) +
				0.114*float64(c.B)

			grayPixels[y*width+x] = gray
		}
	}

	var hash uint64

	for y := 0; y < height; y++ {
		for x := 0; x < width-1; x++ {
			idx := y*width + x
			hash <<= 1

			if grayPixels[idx] > grayPixels[idx+1] {
				hash |= 1
			}
		}
	}

	return hash
}

func hashSimilarityFast(hash1, hash2 uint64) float64 {
	xor := hash1 ^ hash2

	distance := hammingTable[byte(xor>>56)] +
		hammingTable[byte(xor>>48)] +
		hammingTable[byte(xor>>40)] +
		hammingTable[byte(xor>>32)] +
		hammingTable[byte(xor>>24)] +
		hammingTable[byte(xor>>16)] +
		hammingTable[byte(xor>>8)] +
		hammingTable[byte(xor)]

	return 1 - float64(distance)/64.0
}

func colorHistogramSimilarity(img1, img2 *image.RGBA) float64 {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx()*bounds1.Dy() == 0 ||
		bounds2.Dx()*bounds2.Dy() == 0 {
		return 0
	}

	const bins = 8
	hist1 := make([][][]int, bins)
	hist2 := make([][][]int, bins)

	for i := 0; i < bins; i++ {
		hist1[i] = make([][]int, bins)
		hist2[i] = make([][]int, bins)
		for j := 0; j < bins; j++ {
			hist1[i][j] = make([]int, bins)
			hist2[i][j] = make([]int, bins)
		}
	}

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			c := img1.RGBAAt(x, y)
			rIdx := int(c.R) * bins / 256
			gIdx := int(c.G) * bins / 256
			bIdx := int(c.B) * bins / 256
			if rIdx >= bins {
				rIdx = bins - 1
			}
			if gIdx >= bins {
				gIdx = bins - 1
			}
			if bIdx >= bins {
				bIdx = bins - 1
			}
			hist1[rIdx][gIdx][bIdx]++
		}
	}

	for y := bounds2.Min.Y; y < bounds2.Max.Y; y++ {
		for x := bounds2.Min.X; x < bounds2.Max.X; x++ {
			c := img2.RGBAAt(x, y)
			rIdx := int(c.R) * bins / 256
			gIdx := int(c.G) * bins / 256
			bIdx := int(c.B) * bins / 256
			if rIdx >= bins {
				rIdx = bins - 1
			}
			if gIdx >= bins {
				gIdx = bins - 1
			}
			if bIdx >= bins {
				bIdx = bins - 1
			}
			hist2[rIdx][gIdx][bIdx]++
		}
	}

	var similarity float64
	for i := 0; i < bins; i++ {
		for j := 0; j < bins; j++ {
			for k := 0; k < bins; k++ {
				h1 := float64(hist1[i][j][k])
				h2 := float64(hist2[i][j][k])

				if h1 > 0 && h2 > 0 {
					similarity += math.Sqrt(h1 * h2)
				}
			}
		}
	}

	total1 := float64(bounds1.Dx() * bounds1.Dy())
	total2 := float64(bounds2.Dx() * bounds2.Dy())

	if total1 == 0 || total2 == 0 {
		return 0
	}

	return similarity / math.Sqrt(total1*total2)
}

func ssim(img1, img2 *image.RGBA) float64 {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if !bounds1.Eq(bounds2) {
		return 0
	}

	width := bounds1.Dx()
	height := bounds1.Dy()

	if width > 128 || height > 128 {
		return fastSsim(img1, img2)
	}

	if width == 0 || height == 0 {
		return 0
	}

	var sum1, sum2, sum1Sq, sum2Sq, sum12 float64

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c1 := img1.RGBAAt(x+bounds1.Min.X, y+bounds1.Min.Y)
			c2 := img2.RGBAAt(x+bounds2.Min.X, y+bounds2.Min.Y)

			gray1 := 0.299*float64(c1.R) + 0.587*float64(c1.G) + 0.114*float64(c1.B)
			gray2 := 0.299*float64(c2.R) + 0.587*float64(c2.G) + 0.114*float64(c2.B)

			sum1 += gray1
			sum2 += gray2
			sum1Sq += gray1 * gray1
			sum2Sq += gray2 * gray2
			sum12 += gray1 * gray2
		}
	}

	n := float64(width * height)

	mean1 := sum1 / n
	mean2 := sum2 / n
	var1 := (sum1Sq / n) - (mean1 * mean1)
	var2 := (sum2Sq / n) - (mean2 * mean2)
	cov := (sum12 / n) - (mean1 * mean2)

	const C1 = 6.5025  // (0.01*255)^2
	const C2 = 58.5225 // (0.03*255)^2

	numerator := (2*mean1*mean2 + C1) * (2*cov + C2)
	denominator := (mean1*mean1 + mean2*mean2 + C1) * (var1 + var2 + C2)

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

func fastSsim(img1, img2 *image.RGBA) float64 {
	const targetSize = 64

	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if !bounds1.Eq(bounds2) {
		return 0
	}

	width := bounds1.Dx()
	height := bounds1.Dy()

	scale := math.Max(float64(width)/targetSize, float64(height)/targetSize)
	if scale <= 1 {
		return ssim(img1, img2)
	}

	newWidth := int(float64(width) / scale)
	newHeight := int(float64(height) / scale)

	if newWidth < 2 || newHeight < 2 {
		newWidth = 2
		newHeight = 2
	}

	var sum1, sum2, sum1Sq, sum2Sq, sum12 float64

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) * scale)
			srcY := int(float64(y) * scale)

			if srcX >= width {
				srcX = width - 1
			}
			if srcY >= height {
				srcY = height - 1
			}

			c1 := img1.RGBAAt(srcX+bounds1.Min.X, srcY+bounds1.Min.Y)
			c2 := img2.RGBAAt(srcX+bounds2.Min.X, srcY+bounds2.Min.Y)

			gray1 := 0.299*float64(c1.R) + 0.587*float64(c1.G) + 0.114*float64(c1.B)
			gray2 := 0.299*float64(c2.R) + 0.587*float64(c2.G) + 0.114*float64(c2.B)

			sum1 += gray1
			sum2 += gray2
			sum1Sq += gray1 * gray1
			sum2Sq += gray2 * gray2
			sum12 += gray1 * gray2
		}
	}

	n := float64(newWidth * newHeight)

	mean1 := sum1 / n
	mean2 := sum2 / n
	var1 := (sum1Sq / n) - (mean1 * mean1)
	var2 := (sum2Sq / n) - (mean2 * mean2)
	cov := (sum12 / n) - (mean1 * mean2)

	const C1 = 6.5025  // (0.01*255)^2
	const C2 = 58.5225 // (0.03*255)^2

	numerator := (2*mean1*mean2 + C1) * (2*cov + C2)
	denominator := (mean1*mean1 + mean2*mean2 + C1) * (var1 + var2 + C2)

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
