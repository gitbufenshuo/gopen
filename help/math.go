package help

import "strconv"

func Str2Float32(input string) float32 {
	n, _ := strconv.ParseFloat(input, 64)
	return float32(n)
}
