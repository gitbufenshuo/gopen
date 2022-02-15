package logic_jump

import (
	"fmt"
	"math"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
)

// 	PlayerMode_Static   状态下的一些逻辑

// 进入 状态
func (lj *LogicJump) EnterStateStatic() {
	lj.PlayerMode = PlayerMode_Static
	lj.changeACMode("__init") // 切换动画
}

// 离开 状态
func (lj *LogicJump) LeaveStateStatic() {
}

// 状态 步进
func (lj *LogicJump) OnStaticUpdate() {
}

func (lj *LogicJump) OnStaticProcessMSG_move(msg *jump.JumpMSGOne) {
	lj.movex = msg.MoveValX
	lj.movez = msg.MoveValZ
	lj.EnterStateMoving()
}

func (lj *LogicJump) OnStaticProcessMSG_doatt(msg *jump.JumpMSGOne) {
	lj.EnterStateDoAtt()
}

func (lj *LogicJump) OnStaticProcessMSG_underatt(msg *jump.JumpMSGOne) {
	// 静止状态受到了攻击
	underx := lj.logicposx - msg.PosX
	underz := lj.logicposz - msg.PosZ
	mo := math.Sqrt(float64((underx*underx + underz*underz)))
	moint64 := int64(mo)
	if moint64 == 0 {
		moint64 = 1
	}
	underx = underx * 5000 / moint64
	underz = underz * 5000 / moint64
	fmt.Println("static underatt", underx, underz)
	lj.underattx = underx
	lj.underattz = underz
	lj.LeaveStateStatic()   // 退出 static
	lj.EnterStateUnderAtt() // 进入 underatt
}

// 状态 接收msg
func (lj *LogicJump) OnStaticProcessMSG(msg *jump.JumpMSGOne) {
	if msg.Kind == "move" {
		lj.OnStaticProcessMSG_move(msg)
		return
	}
	if msg.Kind == "doatt" {
		lj.OnStaticProcessMSG_doatt(msg)
		return
	}
	if msg.Kind == "underatt" {
		lj.OnStaticProcessMSG_underatt(msg)
		return
	}
}
