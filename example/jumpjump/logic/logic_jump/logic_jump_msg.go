package logic_jump

import "github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"

func (lj *LogicJump) ProcessMSG_move(msg *jump.JumpMSGOne) {
	lj.Velx, lj.Velz = msg.MoveValX, msg.MoveValZ
	lj.Velx = (lj.Velx * lj.moveSpeed) / 100
	lj.Velz = (lj.Velz * lj.moveSpeed) / 100
}

func (lj *LogicJump) ProcessMSG_doatt(msg *jump.JumpMSGOne) {

}
func (lj *LogicJump) ProcessMSG_skill_yingfenshen(msg *jump.JumpMSGOne) {

}

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

	if msg.Kind == "move" {
		lj.ProcessMSG_move(msg)
		return
	}
	if msg.Kind == "doatt" {
		lj.ProcessMSG_doatt(msg)
		return
	}
	if msg.Kind == "skill-yingfenshen" {
		lj.ProcessMSG_skill_yingfenshen(msg)
		return
	}
}
