package stmyplane

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type MyPlane struct {
	*gameobjects.PlaneObject
	CursorMode    int
	KeyEscapeMode int // 0 for release
}

func (myplane *MyPlane) Update() {
	if true {
		if myplane.KeyEscapeMode == 1 {
			myplane.KeyEscapeMode = 0
			if myplane.GI().CursorMode == glfw.CursorNormal {
				myplane.GI().SetCursorMode(glfw.CursorDisabled)
			} else {
				myplane.GI().SetCursorMode(glfw.CursorNormal)
			}
		}
	} else {
		myplane.KeyEscapeMode = 1
	}
	myplane.mouseLook()
}

func (myplane *MyPlane) mouseLook() {
	// gi := myplane.GI()
	// fmt.Printf("mouse diff , %f %f\n", gi.MouseXDiff, gi.MouseYDiff)
	// mainCamera := myplane.GI().MainCamera
	// mainCamera.RotateLocalHorizontal(-float32(gi.MouseXDiff) / 60)
	// mainCamera.RotateLocalVertical(-float32(gi.MouseYDiff) / 60)
}

func NewMyPlane(gi *game.GlobalInfo, modelname, texturename string) *MyPlane {
	myplane := new(MyPlane)
	myplane.PlaneObject = gameobjects.NewPlane(gi, modelname, texturename)
	myplane.KeyEscapeMode = 0
	return myplane
}
