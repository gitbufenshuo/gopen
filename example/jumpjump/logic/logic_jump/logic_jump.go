package logic_jump

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type PlayerMode int

const (
	PlayerMode_Static PlayerMode = 1 // 静止
	PlayerMode_Energy PlayerMode = 2 // 蓄力(动量)
	PlayerMode_Jump   PlayerMode = 3 // 跳
)

type LogicJump struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	//
	playerMode           PlayerMode
	beginms              float64
	velx, vely           float32 // 当前速度
	logicposx, logicposy float32
	gravity              float32
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

func (lj *LogicJump) Start(gb game.GameObjectI) {
	inputsystem.InitInputSystem(lj.gi)
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeySpace))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyA))
	inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyD))
	lj.gi.SetInputSystem(inputsystem.GetInputSystem())
	lj.ac = lj.gi.AnimationSystem.GetAC(gb.ID_sg())
}

func (lj *LogicJump) Update(gb game.GameObjectI) {
	lj.onForce(gb)
	lj.syncLogicPosY(gb)

	if lj.playerMode == PlayerMode_Static {
		lj.PlayerMode_StaticUpdate(gb)
		return
	}
	if lj.playerMode == PlayerMode_Energy {
		lj.PlayerMode_EnergyUpdate(gb)
		return
	}
	if lj.playerMode == PlayerMode_Jump {
		lj.PlayerMode_JumpUpdate(gb)
		return
	}
}

func (lj *LogicJump) onAD(gb game.GameObjectI) {
	apressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyA))
	dpressed := inputsystem.GetInputSystem().KeyPress(int(glfw.KeyD))
	nowmode := lj.ac.NowMode()
	if apressed {
		lj.velx = -5
		if nowmode != "MOVING" {
			lj.ac.ChangeMode("MOVING")
		}
		return
	}
	if dpressed {
		lj.velx = 5
		if nowmode != "MOVING" {
			lj.ac.ChangeMode("MOVING")
		}
		return
	}
	if nowmode == "MOVING" {
		lj.ac.ChangeMode("__init")
	}
	lj.velx = 0
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
	lj.onAD(gb)
	lj.logicposx += lj.velx * deltams
	// fmt.Printf("lj.logicposy:%f lj.vel:%f imp:%f mode:%v\n", lj.logicposy, lj.vel, mergeforce*deltams, lj.playerMode)
}

func (lj *LogicJump) syncLogicPosY(gb game.GameObjectI) {
	gb.GetTransform().Postion.SetValue2(
		lj.logicposx, lj.logicposy,
	)
}

func (lj *LogicJump) PlayerMode_StaticUpdate(gb game.GameObjectI) {
	///////////////////////
	if inputsystem.GetInputSystem().KeyDown(int(glfw.KeySpace)) {
		lj.playerMode = PlayerMode_Jump
		lj.logicposy = 0.0001
		lj.vely = 30
		lj.ac.ChangeMode("MOVING")
	}
}
func (lj *LogicJump) PlayerMode_EnergyUpdate(gb game.GameObjectI) {

}
func (lj *LogicJump) PlayerMode_JumpUpdate(gb game.GameObjectI) {
	fmt.Println("PlayerMode_JumpUpdate")
	if lj.logicposy < 0 {
		lj.playerMode = PlayerMode_Static
		lj.ac.ChangeMode("__init")
	}
	// deltams := lj.gi.FrameElapsedMS
	// lj.energy
}

func (lj *LogicJump) Clone() game.LogicSupportI {
	return NewLogicJump(lj.gi)
}
