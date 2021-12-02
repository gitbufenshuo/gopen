package main

import (
	"os"

	"github.com/gitbufenshuo/gopen/matmath"
)

func TryLookAt() {
	pos := matmath.GetVECX(3)
	pos.SetValue3(0, 0, 1)
	target := matmath.GetVECX(3)
	target.SetValue3(1, 0, 0)
	up := matmath.GetVECX(3)
	up.SetValue3(0, 1, 0)
	viewMAT := matmath.LookAtFrom4(pos, target, up)
	//////
	objectPos := matmath.GetVECX(4)
	objectPos.SetValue4(1, 0, 0, 1)
	//////
	objectPos.RightMul_InPlace(&viewMAT)

	objectPos.PrettyShow()
}

func main() {
	if os.Args[1] == "lookat" {
		TryLookAt()
		return
	}
}
