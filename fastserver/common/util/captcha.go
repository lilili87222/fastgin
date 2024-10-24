package util

import (
	"math"
	"math/rand"
)

const AngleSpin = 20 //角度偏差
func EqualCaptcha(input, expect float64) bool {
	if math.Abs(float64(int(input+expect)%360)) > AngleSpin {
		return false
	}
	return true
}
func RandCaptchaAngle() float64 {
	return AngleSpin + rand.Float64()*(340-AngleSpin)
}
