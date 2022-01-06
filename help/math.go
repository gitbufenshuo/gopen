package help

import (
	"fmt"
	"math"
	"strconv"
)

func Str2Float32(input string) float32 {
	n, _ := strconv.ParseFloat(input, 64)
	return float32(n)
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

func Sin(input float32) float32 {
	return float32(math.Sin(float64(input)))
}
func Cos(input float32) float32 {
	return float32(math.Cos(float64(input)))
}
func ArcCos(input float32) float32 {
	fmt.Println("ArcCos", input)
	return float32(math.Acos(float64(input)))
}
func Sqrt(input float32) float32 {
	return float32(math.Sqrt(float64(input)))
}
