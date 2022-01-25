package logic_jump

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/help"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type PlayerMode int

const (
	PlayerMode_Static PlayerMode = 1 // 静止
	PlayerMode_Jump   PlayerMode = 2 // 跳
)

type LogicJump struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	//
	playerMode                      PlayerMode
	Chosen                          bool
	beginms                         float64
	velx, vely, velz                float32 // 当前速度
	logicposx, logicposy, logicposz float32
	logicroty                       float32
	gravity                         float32
	//
	ac game.AnimationControllerI
}

func NewLogicJump(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicJump)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	res.playerMode = PlayerMode_Jump
	res.gravity = -10 //
	res.logicposx, res.logicposy = 0, 30
	return res
}

func (lj *LogicJump) GetLogicPosX() float32 {
	return lj.logicposx
}
func (lj *LogicJump) getAC(gb game.GameObjectI) {
	if lj.ac != nil {
		return
	}
	lj.ac = lj.gi.AnimationSystem.GetAC(gb.ID_sg())
}
func (lj *LogicJump) changeACMode(mode string) {
	if lj.ac == nil {
		return
	}
	lj.ac.ChangeMode("MOVING")
}

func (lj *LogicJump) Start(gb game.GameObjectI) {
	fmt.Println("logic_jump START")
	inputsystem.InitInputSystem(lj.gi)
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeySpace))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyA))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyD))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyW))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyS))
	lj.gi.SetInputSystem(inputsystem.GetInputSystem())
}

func (lj *LogicJump) Update(gb game.GameObjectI) {
	lj.getAC(gb)
	lj.onAD(gb)
	lj.onForce(gb)
	lj.syncLogicPosY(gb)

	if lj.playerMode == PlayerMode_Static {
		lj.PlayerMode_StaticUpdate(gb)
		return
	}
	if lj.playerMode == PlayerMode_Jump {
		lj.PlayerMode_JumpUpdate(gb)
		return
	}
}

func (lj *LogicJump) onAD(gb game.GameObjectI) {
	lj.velx = 0
	lj.velz = 0
	if !lj.Chosen {
		return
	}
	apressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyA))
	dpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyD))
	wpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyW))
	spressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyS))
	nowmode := lj.ac.NowMode()
	var moved bool
	if apressed {
		lj.velx = -5
		moved = true
	} else if dpressed {
		lj.velx = 5
		moved = true
	}
	if wpressed {
		lj.velz = -5
		moved = true
	} else if spressed {
		lj.velz = 5
		moved = true
	}
	if moved {
		if nowmode != "MOVING" {
			lj.changeACMode("MOVING")
		}
		mo := help.Sqrt(lj.velx*lj.velx+lj.velz*lj.velz) / 5
		lj.velx /= mo
		lj.velz /= mo
		if lj.velx == 0 {
			if lj.velz > 0 {
				lj.logicroty = 0
			} else {
				lj.logicroty = 180
			}
			return
		} else if lj.velx > 0 {
			if lj.velz > 0 {
				lj.logicroty = 45
			} else if lj.velz < 0 {
				lj.logicroty = 135
			} else {
				lj.logicroty = 90
			}
		} else {
			if lj.velz > 0 {
				lj.logicroty = -45
			} else if lj.velz < 0 {
				lj.logicroty = -135
			} else {
				lj.logicroty = -90
			}
		}
		return
	}
	if nowmode == "MOVING" {
		lj.changeACMode("__init")
	}
	lj.velx = 0
	lj.velz = 0
}

func (lj *LogicJump) onForce(gb game.GameObjectI) {
	var upForce float32
	if lj.logicposy <= 0 {
		lj.logicposy = 0
		upForce = -lj.gravity // 如果在地面，向上的弹力应该正好与重力相反
		lj.vely = 0
	}
	//
	deltams := float32(lj.gi.FrameElapsedMS / 1000) // 单位变成秒
	mergeforce := lj.gravity + upForce              // 合力
	lj.vely += (mergeforce * deltams) * 10
	lj.logicposy += lj.vely * deltams
	lj.logicposx += lj.velx * deltams
	lj.logicposz += lj.velz * deltams
	// fmt.Printf("lj.logicposy:%f lj.vel:%f imp:%f mode:%v\n", lj.logicposy, lj.vel, mergeforce*deltams, lj.playerMode)
}

func (lj *LogicJump) syncLogicPosY(gb game.GameObjectI) {
	gb.GetTransform().Postion.SetValue3(
		lj.logicposx, lj.logicposy, lj.logicposz,
	)
	rawroty := gb.GetTransform().Rotation.GetIndexValue(1)
	gb.GetTransform().Rotation.SetIndexValue(1, (lj.logicroty-rawroty)/10+rawroty)
}

func (lj *LogicJump) PlayerMode_StaticUpdate(gb game.GameObjectI) {
	///////////////////////
	if !lj.Chosen {
		return
	}
	if inputsystem.GetInputSystem().KeyDown(int(glfw.KeySpace)) {
		lj.playerMode = PlayerMode_Jump
		lj.logicposy = 0.0001
		lj.vely = 30
		lj.changeACMode("MOVING")
	}
}

func (lj *LogicJump) PlayerMode_JumpUpdate(gb game.GameObjectI) {
	if lj.logicposy < 0 {
		lj.playerMode = PlayerMode_Static
		lj.changeACMode("__init")
	}
	// deltams := lj.gi.FrameElapsedMS
	// lj.energy
}

func (lj *LogicJump) Clone() game.LogicSupportI {
	return NewLogicJump(lj.gi)
}
