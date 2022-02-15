package logic_jump

import (
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/share/pkem"
)

// 	PlayerMode_UnderAtt   状态下的一些逻辑

// 进入 状态
func (lj *LogicJump) EnterStateDoAtt() {
	lj.PlayerMode = PlayerMode_DoAtt
	lj.movex = 0
	lj.movez = 0
	//
	lj.doattmsFrame = lj.outterFrame
	//
	lj.changeACMode("pugong")
}

// 离开 状态
func (lj *LogicJump) LeaveStateDoAtt() {
}

// 状态 步进
func (lj *LogicJump) OnDoAttUpdate() {
	if lj.outterFrame-lj.doattmsFrame > 8 {
		// 此时释放攻击,并切换到static状态
		lj.LeaveStateDoAtt()
		lj.EnterStateStatic()
		// fire event
		{
			newev := new(pkem.Event)
			// PID              int64     // playerid
			// EK               EventKind // 事件类型
			// PosX, PosY, PosZ int64     // 事件发生的坐标
			// DirX, DirY, DirZ int64     // 事件发生的朝向
			newev.PID = lj.pid
			newev.EK = pkem.EK_DoAtt
			newev.PosX = lj.logicposx
			newev.PosY = lj.logicposy
			newev.PosZ = lj.logicposz
			lj.evm.FireEvent(newev)
		}
	}
}
func (lj *LogicJump) OnDoAttProcessMSG_move(msg *jump.JumpMSGOne) {
	lj.movex = msg.MoveValX
	lj.movez = msg.MoveValZ
	lj.LeaveStateDoAtt()
	lj.EnterStateMoving()
}

func (lj *LogicJump) OnDoAttProcessMSG_doatt(msg *jump.JumpMSGOne) {
	// 忽略
}

// 状态 接收msg
func (lj *LogicJump) OnDoAttProcessMSG(msg *jump.JumpMSGOne) {
	if msg.Kind == "move" {
		lj.OnDoAttProcessMSG_move(msg)
		return
	}
	if msg.Kind == "doatt" {
		lj.OnDoAttProcessMSG_doatt(msg)
		return
	}
}
