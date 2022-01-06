package main

import (
	"fmt"
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

func Angle() {
	bound1 := matmath.CreateVec2(0, 0)
	bound2 := matmath.CreateVec2(1, 0)
	bound3 := matmath.CreateVec2(1, 1)
	bound4 := matmath.CreateVec2(0, 1)
	target := matmath.CreateVec2(0.5, 0.5)
	//
	boundlist := make([]*matmath.Vec2, 4)
	boundlist[0] = &bound1
	boundlist[1] = &bound2
	boundlist[2] = &bound3
	boundlist[3] = &bound4
	angle := matmath.Vec2BoundCheck(boundlist, &target)
	fmt.Println(angle)
}

func main() {
	if os.Args[1] == "lookat" {
		TryLookAt()
		return
	}
	if os.Args[1] == "angle" {
		Angle()
		return
	}
}
