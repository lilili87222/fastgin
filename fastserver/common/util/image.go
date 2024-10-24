package util

import (
	"bytes"
	"encoding/base64"
	_ "golang.org/x/image/bmp"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"math"
	"os"
)

func RotateImage(img image.Image, angle float64) (image.Image, error) {
	// 计算旋转后的图像的尺寸
	radians := angle * math.Pi / 180
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	newWidth := int(math.Abs(float64(width*int(math.Cos(radians))) + float64(height*int(math.Sin(radians)))))
	newHeight := int(math.Abs(float64(width*int(math.Sin(radians))) + float64(height*int(math.Cos(radians)))))

	// 创建一个新的图像
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 计算旋转中心
	centerX, centerY := float64(width)/2, float64(height)/2
	newCenterX, newCenterY := float64(newWidth)/2, float64(newHeight)/2

	// 旋转图像
	for x := 0; x < newWidth; x++ {
		for y := 0; y < newHeight; y++ {
			// 反向计算原图中的坐标
			originalX := int(centerX + (float64(x)-newCenterX)*math.Cos(-radians) - (float64(y)-newCenterY)*math.Sin(-radians))
			originalY := int(centerY + (float64(x)-newCenterX)*math.Sin(-radians) + (float64(y)-newCenterY)*math.Cos(-radians))

			// 判断原坐标是否在原图范围内
			if originalX >= 0 && originalX < width && originalY >= 0 && originalY < height {
				newImg.Set(x, y, img.At(originalX, originalY))
			}
		}
	}

	return newImg, nil
}

// 扭曲函数
func DistortImage(img image.Image, strength float64) (image.Image, error) {
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

	return newImage, nil
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
	img, err = DistortImage(img, 20)
	if err != nil {
		return "", err
	}
	// Rotate the image
	rotatedImg, err := RotateImage(img, angle)
	if err != nil {
		return "", err
	}

	// Encode the rotated image to JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rotatedImg, nil)
	if err != nil {
		return "", err
	}

	// Convert the JPEG to base64
	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Img, nil
}
