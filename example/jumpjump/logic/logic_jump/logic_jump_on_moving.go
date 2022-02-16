package logic_jump

import (
	"fmt"
	"math"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
)

// 	PlayerMode_Moving   状态下的一些逻辑

// 进入 状态
func (lj *LogicJump) EnterStateMoving() {
	lj.PlayerMode = PlayerMode_Moving
	lj.changeACMode("yidong") // 切换动画
}

// 离开 状态
func (lj *LogicJump) LeaveStateMoving() {
}

// 状态 步进
func (lj *LogicJump) OnMovingUpdate() {
	lj.movex /= 3
	lj.movez /= 3
	if lj.movex == 0 && lj.movez == 0 {
		lj.LeaveStateMoving()
		lj.EnterStateStatic()
	}
}
func (lj *LogicJump) OnMovingProcessMSG_move(msg *jump.JumpMSGOne) {
	lj.movex = msg.MoveValX
	lj.movez = msg.MoveValZ
}

func (lj *LogicJump) OnMovingProcessMSG_doatt(msg *jump.JumpMSGOne) {
	lj.EnterStateDoAtt()
}

func (lj *LogicJump) OnMovingProcessMSG_underatt(msg *jump.JumpMSGOne) {
	// 移动状态 状态受到了攻击
	underx := lj.logicposx - msg.PosX
	underz := lj.logicposz - msg.PosZ
	mo := math.Sqrt(float64((underx*underx + underz*underz)))
	moint64 := int64(mo)
	if moint64 == 0 {
		moint64 = 1
	}
	if moint64 > 2000 {
		return
	}
	underx = underx * 5000 / moint64
	underz = underz * 5000 / moint64
	fmt.Println("moving underatt", underx, underz)
	lj.underattx = underx
	lj.underattz = underz
	lj.LeaveStateMoving()   // 退出 static
	lj.EnterStateUnderAtt() // 进入 underatt
}

// 状态 接收msg
func (lj *LogicJump) OnMovingProcessMSG(msg *jump.JumpMSGOne) {
	if msg.Kind == "move" {
		lj.OnMovingProcessMSG_move(msg)
		return
	}
	if msg.Kind == "doatt" {
		lj.OnMovingProcessMSG_doatt(msg)
		return
	}
	if msg.Kind == "underatt" {
		lj.OnStaticProcessMSG_underatt(msg)
		return
	}

}
