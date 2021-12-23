package help

import "strconv"

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
