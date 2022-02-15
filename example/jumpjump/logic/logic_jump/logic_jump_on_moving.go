package logic_jump

import "github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"

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
	lj.movex /= 2
	lj.movez /= 2
	if lj.movex == 0 && lj.movez == 0 {
		lj.LeaveStateMoving()
		lj.EnterStateStatic()
	}
}
func (lj *LogicJump) OnMovingProcessMSG_move(msg *jump.JumpMSGOne) {
	lj.movex = msg.MoveValX
	lj.movez = msg.MoveValZ
}

// 状态 接收msg
func (lj *LogicJump) OnMovingProcessMSG(msg *jump.JumpMSGOne) {
	if msg.Kind == "move" {
		lj.OnMovingProcessMSG_move(msg)
		return
	}

}
