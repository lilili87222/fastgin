package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

const AngleSpin = 20 //角度偏差
func EqualCaptcha(input, myAngle float64) bool {
	expect := 360 - myAngle
	if math.Abs(input-expect) > AngleSpin {
		return false
	}
	return true
}
func RandCaptchaAngle() float64 {
	return 5 + rand.Float64()*350
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
	img = DrawLines(img, 15, color.White, 2)

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
