package util

import (
	"math"
	"math/rand"
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
