package logic_jump

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/example/htmlblockcustom/logic/logic_follow"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type PlayerMode int

const (
	PlayerMode_Static   PlayerMode = 1 // 静止
	PlayerMode_Jump     PlayerMode = 2 // 跳
	PlayerMode_UnderAtt PlayerMode = 3 // 受攻击
	PlayerMode_DoAtt    PlayerMode = 4 // 主动攻击
)

type LogicJump struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	Transform *game.Transform
	//
	PlayerMode                      PlayerMode
	Chosen                          bool
	beginms                         float64
	Velx, Vely, Velz                int64 // 当前速度
	logicposx, logicposy, logicposz int64
	rlogicposx, rlogicposz          int64
	Logicroty                       int64 // 1 代表 0.01°
	factor                          float32
	gravity                         int64
	frame                           int
	ljs                             *LogicJumpSignal
	//
	fenshenList []*logic_follow.LogicFollow
	ac          game.AnimationControllerI
}

type LogicJumpSignal struct {
	Kind     string // move jump
	MoveValX int64
	MoveValZ int64
}

func NewLogicJump(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicJump)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	res.PlayerMode = PlayerMode_Static
	res.gravity = -10 //
	res.logicposx, res.logicposy = 0, 0
	res.ljs = new(LogicJumpSignal)
	res.factor = 5
	return res
}

func (lj *LogicJump) GetLogicPosX() int64 {
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
	lj.Transform = gb.GetTransform()
}

func (lj *LogicJump) Update(gb game.GameObjectI) {
	lj.frame++
	lj.getAC(gb)         // 逻辑无关
	lj.syncLogicPosY(gb) // 逻辑无关
	return
	if lj.PlayerMode == PlayerMode_Static {
		lj.PlayerMode_StaticUpdate(gb)
		return
	}
	if lj.PlayerMode == PlayerMode_Jump {
		lj.PlayerMode_JumpUpdate(gb)
		return
	}
}

func (lj *LogicJump) OutterUpdate() {
	lj.OnForce()
	for idx, onefenshen := range lj.fenshenList {
		onefenshen.Move(lj.logicposx, lj.logicposy, lj.logicposz, int64(idx*5))
	}
	//
	if lj.PlayerMode == PlayerMode_UnderAtt {
		lj.Velx = help.Int64To(lj.Velx, 0, 90)
		lj.Velz = help.Int64To(lj.Velz, 0, 90)
		if lj.Velx == 0 && lj.Velz == 0 {
			lj.PlayerMode = PlayerMode_Static
		}
		return
	}
	if lj.PlayerMode == PlayerMode_DoAtt {
		lj.logicposx, lj.logicposz = help.Int64To(lj.logicposx, lj.rlogicposx, 80), help.Int64To(lj.logicposz, lj.rlogicposz, 80)
		if lj.logicposx == lj.rlogicposx && lj.logicposz == lj.rlogicposz {
			lj.PlayerMode = PlayerMode_Static
		}
		return
	}
}

func (lj *LogicJump) EnterPlayerMode_DoAtt() {
	lj.PlayerMode = PlayerMode_DoAtt
	lj.Velx = 0
	lj.Velz = 0
	lj.rlogicposx, lj.rlogicposz = lj.logicposx, lj.logicposz
	lj.logicposx += 3000
	lj.logicposz += 3000
}

func (lj *LogicJump) Skill_Yingfenshen() {
	// instantiate dunshan1
	prefab := modelcustom.PrefabSystemIns.GetPrefab("dunshanying")
	if prefab != nil {
		gb := prefab.Instantiate(lj.gi)
		gb.GetTransform().Postion.SetIndexValue(0, lj.Transform.Postion.GetIndexValue(0))
		logiclist := gb.GetLogicSupport()
		for idx := range logiclist {
			if v, ok := logiclist[idx].(*logic_follow.LogicFollow); ok {
				lj.fenshenList = append(lj.fenshenList, v)
			}
		}
	}
}

func (lj *LogicJump) EnterPlayerMode_UnderAtt() {
	lj.PlayerMode = PlayerMode_UnderAtt
	forward := lj.Transform.GetForward()
	fx, _, fz := forward.GetValue3()
	fmt.Println("forward", fx, fz)
	lj.Velx = int64(fx * 5000)
	lj.Velz = int64(fz * 5000)
}

func (lj *LogicJump) OnForce() {
	var upForce int64
	if lj.logicposy <= 0 {
		lj.logicposy = 0
		upForce = -lj.gravity // 如果在地面，向上的弹力应该正好与重力相反
		lj.Vely = 0
	}
	//
	//deltams := float32(lj.gi.FrameElapsedMS / 1000) // 单位变成秒
	mergeforce := lj.gravity + upForce // 合力
	lj.Vely += (mergeforce) * 10
	lj.logicposy += lj.Vely
	lj.logicposx += lj.Velx
	lj.logicposz += lj.Velz
	lj.factor = 5
	{
		// clamp x and z
		if lj.logicposx < -20*1000 {
			lj.logicposx = -20 * 1000
		}
		if lj.logicposx > 20*1000 {
			lj.logicposx = 20 * 1000
		}

		if lj.logicposz < -10*1000 {
			lj.logicposz = -10 * 1000
		}
		if lj.logicposz > 10*1000 {
			lj.logicposz = 10 * 1000
		}
	}

	// fmt.Printf("lj.logicposy:%f lj.vel:%f imp:%f mode:%v\n", lj.logicposy, lj.vel, mergeforce*deltams, lj.PlayerMode)
}

func (lj *LogicJump) syncLogicPosY(gb game.GameObjectI) {
	nowposx, nowposy, nowposz := gb.GetTransform().Postion.GetValue3()
	nowposx += (float32(lj.logicposx)/1000 - nowposx) / lj.factor
	nowposy += (float32(lj.logicposy)/1000 - nowposy) / lj.factor
	nowposz += (float32(lj.logicposz)/1000 - nowposz) / lj.factor
	gb.GetTransform().Postion.SetValue3(
		nowposx, nowposy, nowposz,
	)
	forward := matmath.CreateVec4(float32(lj.Velx), 0, float32(lj.Velz), 1)
	gb.GetTransform().SetForward(forward, 0.2)
	// curframe := float32(lj.gi.CurFrame)
	// gb.GetTransform().Rotation.SetIndexValue(0, -30)
	// gb.GetTransform().Rotation.SetIndexValue(1, 45)
	// gb.GetTransform().Rotation.SetIndexValue(1, rawroty)
}

func (lj *LogicJump) PlayerMode_StaticUpdate(gb game.GameObjectI) {
	///////////////////////
	if !lj.Chosen {
		return
	}
	if inputsystem.GetInputSystem().KeyDown(int(glfw.KeySpace)) {
		lj.PlayerMode = PlayerMode_Jump
		lj.logicposy = 1
		lj.Vely = 30
		lj.changeACMode("MOVING")
	}
}

func (lj *LogicJump) PlayerMode_JumpUpdate(gb game.GameObjectI) {
	if lj.logicposy < 0 {
		lj.PlayerMode = PlayerMode_Static
		lj.changeACMode("__init")
	}
	// deltams := lj.gi.FrameElapsedMS
	// lj.energy
}

func (lj *LogicJump) Clone() game.LogicSupportI {
	return NewLogicJump(lj.gi)
}
