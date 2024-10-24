package util

import (
	"bytes"
	"encoding/base64"
	_ "golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

func RotateImage(img image.Image, angle float64) image.Image {
	radians := angle * math.Pi / 180
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	newWidth := int(math.Abs(float64(width)*math.Cos(radians)) + math.Abs(float64(height)*math.Sin(radians)))
	newHeight := int(math.Abs(float64(width)*math.Sin(radians)) + math.Abs(float64(height)*math.Cos(radians)))
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	centerX, centerY := float64(width)/2, float64(height)/2
	newCenterX, newCenterY := float64(newWidth)/2, float64(newHeight)/2

	for x := 0; x < newWidth; x++ {
		for y := 0; y < newHeight; y++ {
			originalX := int(centerX + (float64(x)-newCenterX)*math.Cos(-radians) - (float64(y)-newCenterY)*math.Sin(-radians))
			originalY := int(centerY + (float64(x)-newCenterX)*math.Sin(-radians) + (float64(y)-newCenterY)*math.Cos(-radians))

			if originalX >= 0 && originalX < width && originalY >= 0 && originalY < height {
				newImg.Set(x, y, img.At(originalX, originalY))
			} else {
				newImg.Set(x, y, color.Transparent)
			}
		}
	}
	return newImg
}

// 扭曲函数
func DistortImage(img image.Image, strength float64) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算扭曲后的坐标
			offsetX := int(strength * math.Sin(float64(y)/20))
			offsetY := int(strength * math.Cos(float64(x)/20))
			srcX := x + offsetX
			srcY := y + offsetY

			// 确保坐标在有效范围内
			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				newImage.Set(x, y, img.At(srcX, srcY))
			} else {
				newImage.Set(x, y, color.Transparent)
			}
		}
	}

	return newImage
}
func Blur(img image.Image, radius int) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a uint32
			count := 0

			// 计算周围像素的颜色平均值
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					px := x + dx
					py := y + dy
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						col := img.At(px, py)
						rr, gg, bb, aa := col.RGBA()
						r += rr
						g += gg
						b += bb
						a += aa
						count++
					}
				}
			}

			if count > 0 {
				newImg.Set(x, y, color.RGBA{
					R: uint8(r / uint32(count) >> 8),
					G: uint8(g / uint32(count) >> 8),
					B: uint8(b / uint32(count) >> 8),
					A: uint8(a / uint32(count) >> 8),
				})
			}
		}
	}
	return newImg
}

// DrawLines 在图像上随机绘制指定数量的线条
func DrawLines(img image.Image, count int, col color.Color, width int) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, img, image.Point{}, draw.Src)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		x1 := rand.Intn(bounds.Max.X)
		y1 := rand.Intn(bounds.Max.Y)
		x2 := rand.Intn(bounds.Max.X)
		y2 := rand.Intn(bounds.Max.Y)

		drawLine(newImg, x1, y1, x2, y2, col, width)
	}

	return newImg
}

// drawLine 画线条的函数，支持线宽
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color, width int) {
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	sx := 1
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	sy := 1
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		for w := -width / 2; w <= width/2; w++ {
			// 画出粗线的每一条像素
			if x1+w >= 0 && x1+w < img.Bounds().Max.X && y1 >= 0 && y1 < img.Bounds().Max.Y {
				img.Set(x1+w, y1, col)
			}
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		err2 := err * 2
		if err2 > -dy {
			err -= dy
			x1 += sx
		}
		if err2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// Base64ImageFile reads an image from a file, rotates it by the given angle, and returns the base64 encoded image.
func Base64ImageFile(filePath string, angle float64) (string, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	img = DistortImage(img, 10)

	//img = Blur(img, 2)
	img = DrawLines(img, 20, color.White, 4)

	// Rotate the image
	img = RotateImage(img, angle)

	// Encode the rotated image to JPEG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	//os.WriteFile("test.png", buf.Bytes(), os.ModePerm)

	// Convert the JPEG to base64
	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Img, nil
}
