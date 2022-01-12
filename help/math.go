package help

import (
	"math"
	"strconv"
)

func Str2Float32(input string) float32 {
	n, _ := strconv.ParseFloat(input, 64)
	return float32(n)
}
func Str2Int(input string) int {
	n, _ := strconv.Atoi(input)
	return n
}

func Mi2(input int) int {
	var nowbase int = 2
	for {
		if nowbase > input {
			return nowbase
		}
		nowbase *= 2
	}
}

func TheSame(one, two float32) bool {
	a := one - two
	return math.Abs(float64(a)) < 0.0001
}

func Sin(input float32) float32 {
	return float32(math.Sin(float64(input)))
}
func Cos(input float32) float32 {
	return float32(math.Cos(float64(input)))
}
func ArcCos(input float32) float32 {
	return float32(math.Acos(float64(input)))
}
func Sqrt(input float32) float32 {
	return float32(math.Sqrt(float64(input)))
}
