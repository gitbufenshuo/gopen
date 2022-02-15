package logic_jump

import "github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"

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

// 状态 接收msg
func (lj *LogicJump) OnStaticProcessMSG(msg *jump.JumpMSGOne) {
	if msg.Kind == "move" {
		lj.OnStaticProcessMSG_move(msg)
		return
	}

}
