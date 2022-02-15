package logic_jump

import "github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"

func (lj *LogicJump) ProcessMSG(msg *jump.JumpMSGOne) {
	if lj.PlayerMode == PlayerMode_Static {
		lj.OnStaticProcessMSG(msg)
		return
	}
	if lj.PlayerMode == PlayerMode_Moving {
		lj.OnMovingProcessMSG(msg)
		return
	}
	if lj.PlayerMode == PlayerMode_DoAtt {
		lj.OnDoAttProcessMSG(msg)
		return
	}
}
