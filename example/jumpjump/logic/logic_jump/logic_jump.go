package logic_jump

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/example/htmlblockcustom/logic/logic_follow"
	"github.com/gitbufenshuo/gopen/example/jumpjump/share/pkem"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
	"github.com/gitbufenshuo/gopen/matmath"
)

type PlayerMode int

const (
	PlayerMode_Static   PlayerMode = iota + 1 // 静止
	PlayerMode_Moving   PlayerMode = iota + 1 // 移动
	PlayerMode_Jump     PlayerMode = iota + 1 // 跳
	PlayerMode_UnderAtt PlayerMode = iota + 1 // 受攻击
	PlayerMode_DoAtt    PlayerMode = iota + 1 // 主动攻击
)

type LogicJump struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	Transform *game.Transform
	//
	PlayerMode                      PlayerMode
	pid                             int64
	evm                             *pkem.EventManager
	beginms                         float64
	Velx, Vely, Velz                int64 // 当前总速度
	movex, movez                    int64 // 由于wsad产生的速度
	underattx, underattz            int64 // 由于 被攻击 产生的速度
	logicposx, logicposy, logicposz int64
	rlogicposx, rlogicposz          int64
	Logicroty                       int64 // 1 代表 0.01°
	factor                          float32
	gravity                         int64
	frame                           int
	outterFrame                     int64
	doattmsFrame                    int64
	underattmsFrame                 int64
	moveSpeed                       int64 // 100 代表百分之百
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
	res.factor = 5
	res.moveSpeed = 100
	return res
}

func (lj *LogicJump) GetLogicPosX() int64 {
	return lj.logicposx
}

func (lj *LogicJump) GetLogicPosY() int64 {
	return lj.logicposy
}

func (lj *LogicJump) GetLogicPosZ() int64 {
	return lj.logicposz
}

func (lj *LogicJump) getAC(gb game.GameObjectI) {
	if lj.ac != nil {
		return
	}
	lj.ac = gb.GetACSupport()
}
func (lj *LogicJump) SetPID(pid int64) {
	lj.pid = pid
}
func (lj *LogicJump) GetPID() int64 {
	return lj.pid
}
func (lj *LogicJump) SetEVM(evm *pkem.EventManager) {
	fmt.Println("set evm", evm)
	lj.evm = evm
}

// 切换动画，不影响逻辑
func (lj *LogicJump) changeACMode(mode string) {
	if lj.ac == nil {
		return
	}
	lj.ac.ChangeMode(mode)
}

func (lj *LogicJump) Start(gb game.GameObjectI) {
	fmt.Println("logic_jump START")
	lj.Transform = gb.GetTransform()
}

func (lj *LogicJump) Update(gb game.GameObjectI) {
	// 此为逻辑狗无关的 update 比如说，同步逻辑位置到渲染位置
	lj.frame++
	lj.getAC(gb)         // 逻辑无关
	lj.syncLogicPosY(gb) // 逻辑无关
}

func (lj *LogicJump) OutterUpdate() {
	lj.outterFrame++
	lj.OnForce()
	for idx, onefenshen := range lj.fenshenList {
		onefenshen.Move(lj.logicposx, lj.logicposy, lj.logicposz, int64(idx*5))
	}
	//
	if lj.PlayerMode == PlayerMode_Static {
		lj.OnStaticUpdate()
		return
	}
	if lj.PlayerMode == PlayerMode_Moving {
		lj.OnMovingUpdate()
		return
	}
	if lj.PlayerMode == PlayerMode_DoAtt {
		lj.OnDoAttUpdate()
		return
	}
	if lj.PlayerMode == PlayerMode_UnderAtt {
		lj.OnUnderAttUpdate()
		return
	}
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

func (lj *LogicJump) syncLogicPosY(gb game.GameObjectI) {
	nowposx, nowposy, nowposz := gb.GetTransform().Postion.GetValue3()
	nowposx += (float32(lj.logicposx)/1000 - nowposx) / lj.factor
	nowposy += (float32(lj.logicposy)/1000 - nowposy) / lj.factor
	nowposz += (float32(lj.logicposz)/1000 - nowposz) / lj.factor
	gb.GetTransform().Postion.SetValue3(
		nowposx, nowposy, nowposz,
	)
	forward := matmath.CreateVec4(float32(lj.GetVelX()), 0, float32(lj.GetVelZ()), 1)
	gb.GetTransform().SetForward(forward, 0.2)
	// curframe := float32(lj.gi.CurFrame)
	// gb.GetTransform().Rotation.SetIndexValue(0, -30)
	// gb.GetTransform().Rotation.SetIndexValue(1, 45)
	// gb.GetTransform().Rotation.SetIndexValue(1, rawroty)
}
