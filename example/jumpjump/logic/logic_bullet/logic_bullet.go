package logic_bullet

import (
	"fmt"
	"math"

	"github.com/gitbufenshuo/gopen/example/jumpjump/share/pkem"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

type LogicBullet struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	evm *pkem.EventManager
	//
	TargetPID                          int64 // 目标对象的id
	LogicPosX, LogicPosY, LogicPosZ    int64 // 当前位置
	Vel                                int64 // 当前速度
	TargetPosX, TargetPosY, TargetPosZ int64 // 目标当前位置
	outterframe                        int64
	releaseFrame                       int64
	ShouldDel                          bool
	attTime                            int64 // 攻击次数
}

func NewLogicBullet(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicBullet)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}
func (lj *LogicBullet) SetEVM(evm *pkem.EventManager) {
	fmt.Println("set evm", evm)
	lj.evm = evm
}

func (lj *LogicBullet) Update(gb game.GameObjectI) {
	lj.syncLogicPosY(gb)
}

func (lj *LogicBullet) OutterUpdate() {
	if lj.ShouldDel {
		lj.LogicPosX += 100
		lj.LogicPosZ += 100
		return
	}
	lj.outterframe++
	underx := lj.TargetPosX - lj.LogicPosX
	underz := lj.TargetPosZ - lj.LogicPosZ
	lj.LogicPosX += underx / 10
	lj.LogicPosZ += underz / 10
	mo := math.Sqrt(float64((underx*underx + underz*underz)))
	if mo < 3000 {
		// 释放一个doatt事件
		if lj.outterframe-lj.releaseFrame > 30 {
			lj.releaseFrame = lj.outterframe
			newev := new(pkem.Event)
			// PID              int64     // playerid
			// EK               EventKind // 事件类型
			// PosX, PosY, PosZ int64     // 事件发生的坐标
			// DirX, DirY, DirZ int64     // 事件发生的朝向
			newev.PID = -1
			newev.EK = pkem.EK_DoAtt
			newev.PosX = lj.LogicPosX
			newev.PosY = lj.LogicPosY
			newev.PosZ = lj.LogicPosZ
			{
				newev.TargetPID = lj.TargetPID
			}
			lj.evm.FireEvent(newev)
			fmt.Println("一个小傀儡正在攻击", lj.attTime, newev.TargetPID, newev.EK)
			lj.attTime++
			if lj.attTime > 5 {
				lj.ShouldDel = true
			}
		}
	}
}

func (lj *LogicBullet) syncLogicPosY(gb game.GameObjectI) {
	nowposx, nowposy, nowposz := gb.GetTransform().Postion.GetValue3()
	nowposx += (float32(lj.LogicPosX)/1000 - nowposx) / 2
	nowposy += (float32(lj.LogicPosY)/1000 - nowposy) / 2
	nowposz += (float32(lj.LogicPosZ)/1000 - nowposz) / 2
	gb.GetTransform().Postion.SetValue3(
		nowposx, nowposy, nowposz,
	)
	// forward := matmath.CreateVec4(float32(lj.GetVelX()), 0, float32(lj.GetVelZ()), 1)
	// gb.GetTransform().SetForward(forward, 0.2)
	// curframe := float32(lj.gi.CurFrame)
	// gb.GetTransform().Rotation.SetIndexValue(0, -30)
	// gb.GetTransform().Rotation.SetIndexValue(1, 45)
	// gb.GetTransform().Rotation.SetIndexValue(1, rawroty)
}
