package char_photo

import (
	"github.com/aybabtme/rgbterm"
	_ "github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"os"
	"strings"
)

// 字符集。
var (
	CharSet1 = []byte(`@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\|()1{}[]?-_+~<>i!lI;:,\"^'. `)
	CharSet2 = []byte(`@&W0*+. `)
	CharSet3 = []byte(`@ `)
)

// OpenImage 打开图片，返回图片对象。
func OpenImage(path string) (image.Image, error) {
	// 1. 打开文件。
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	// 2. 解析为图片对象。
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	// 3. 返回。
	return img, nil
}

// PixelMatrix 像素矩阵。
type PixelMatrix struct {
	Matrix [][]*Pixel
}

// ToCharPhotoGray 转为字符图片。 (灰白)
func (m *PixelMatrix) ToCharPhotoGray() string {
	charLines := make([]string, 0)
	for y := 0; y < len(m.Matrix); y++ {
		var charLine string
		for x := 0; x < len(m.Matrix[y]); x++ {
			charLine += string(m.Matrix[y][x].Char)
		}
		charLines = append(charLines, charLine)
	}
	return strings.Join(charLines, "\n")
}

// ToCharPhotoColorful 转为字符图片。 (彩色)
func (m *PixelMatrix) ToCharPhotoColorful() string {
	charLines := make([]string, 0)
	for y := 0; y < len(m.Matrix); y++ {
		var charLine string
		for x := 0; x < len(m.Matrix[y]); x++ {
			pixel := m.Matrix[y][x]
			charLine += DecorateWithColor(pixel.Char, pixel.Color)
		}
		charLines = append(charLines, charLine)
	}
	return strings.Join(charLines, "\n")
}

// ----------------------------------------------

// Pixel 像素。
type Pixel struct {
	Color color.Color
	Char  byte // 字符。
}

// -------------------------------------------------------------------------------------------------------------------

// BuildPixelMatrixByImgPath 构建像素矩阵。
func BuildPixelMatrixByImgPath(path string, width, height uint, charSet []byte) (*PixelMatrix, error) {
	// 1. 打开原始图片。
	orgImg, err := OpenImage(path)
	if err != nil {
		return nil, err
	}

	// 2. 缩放，得到缩放图片。
	scaleImg := resize.Resize(width, height, orgImg, resize.Lanczos3)

	// 3. 形成像素矩阵，返回。
	return BuildPixelMatrix(scaleImg, charSet), nil
}

// BuildPixelMatrix 构建像素矩阵。
func BuildPixelMatrix(img image.Image, charSet []byte) *PixelMatrix {
	matrix := make([][]*Pixel, 0)
	for y := 0; y < img.Bounds().Dy(); y++ { // 方向 ↓。
		pixels := make([]*Pixel, 0)
		for x := 0; x < img.Bounds().Dx(); x++ { // 方向 →。
			// 1. 获取RGB、灰度值。
			pixelColor := img.At(x, y)
			r, g, b, _ := pixelColor.RGBA()
			gray := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)

			// 2. 根据灰度选用字符集的字符。
			charIndex := int(gray / 65536.0 * float64(len(charSet)))

			// 3. 形成像素。
			pixels = append(pixels, &Pixel{
				Color: pixelColor,
				Char:  charSet[charIndex],
			})
		}
		matrix = append(matrix, pixels)
	}
	return &PixelMatrix{Matrix: matrix}
}

// DecorateWithColor 用颜色装饰字符，返回ANSI序列。
// 序列字符串如: "\x1b[38;5;245m9\x1b[0m"，在使用 fmt.Print 输出该字符串，终端会输出一个灰色的 9。
func DecorateWithColor(char byte, c color.Color) string {
	r, g, b, _ := c.RGBA()
	return rgbterm.FgString(string([]byte{char}), uint8(r/256), uint8(g/256), uint8(b/256))
}
