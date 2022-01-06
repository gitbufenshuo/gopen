package main

import (
	"os"

	"github.com/gitbufenshuo/gopen/matmath"
)

func TryLookAt() {
	var pos matmath.Vec4
	pos.SetValue3(0, 0, 1)

	var target matmath.Vec4
	target.SetValue3(1, 0, 0)

	var up matmath.Vec4
	up.SetValue3(0, 1, 0)
	viewMAT := matmath.LookAtFrom4(&pos, &target, &up)
	//////
	var objectPos matmath.Vec4
	objectPos.SetValue4(1, 0, 0, 1)
	//////
	objectPos.RightMul_InPlace(&viewMAT)

	objectPos.PrettyShow("")
}

func main() {
	if os.Args[1] == "lookat" {
		TryLookAt()
		return
	}
}
