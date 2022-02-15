package logic_jump

import "github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"

// 	PlayerMode_Static   状态下的一些逻辑

// 进入 状态
func (lj *LogicJump) EnterStateUnderAtt() {
	lj.PlayerMode = PlayerMode_UnderAtt
}

// 离开 状态
func (lj *LogicJump) LeaveStateUnderAtt() {
	lj.underattx = 0
	lj.underattz = 0
}

// 状态 步进
func (lj *LogicJump) OnUnderAttUpdate() {
	lj.underattx /= 2
	lj.underattz /= 2
	if lj.underattx == 0 && lj.underattz == 0 {
		lj.LeaveStateUnderAtt()
		lj.EnterStateStatic()
	}
}

func (lj *LogicJump) OnUnderAttProcessMSG_move(msg *jump.JumpMSGOne) {

}

func (lj *LogicJump) OnUnderAttProcessMSG_doatt(msg *jump.JumpMSGOne) {

}

func (lj *LogicJump) OnUnderAttProcessMSG_underatt(msg *jump.JumpMSGOne) {

}

// 状态 接收msg
func (lj *LogicJump) OnUnderAttProcessMSG(msg *jump.JumpMSGOne) {
	if msg.Kind == "move" {
		lj.OnUnderAttProcessMSG_move(msg)
		return
	}
	if msg.Kind == "doatt" {
		lj.OnUnderAttProcessMSG_doatt(msg)
		return
	}
	if msg.Kind == "underatt" {
		lj.OnUnderAttProcessMSG_underatt(msg)
		return
	}
}
