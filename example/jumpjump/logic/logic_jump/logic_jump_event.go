package logic_jump

import (
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
	"github.com/gitbufenshuo/gopen/example/jumpjump/share/pkem"
)

func (lj *LogicJump) OutterEV() {
	// 遍历事件
	if lj.evm == nil {
		return
	}
	evlist := lj.evm.GetEvList()
	for _, oneev := range evlist {
		if oneev.EK == pkem.EK_DoAtt { // 如果是攻击事件
			if lj.pid != oneev.PID { // 不是自己
				// 进行攻击判定
				newmsg := new(jump.JumpMSGOne)
				newmsg.Kind = "underatt"
				newmsg.PosX, newmsg.PosY, newmsg.PosZ = oneev.PosX, oneev.PosY, oneev.PosZ
				lj.ProcessMSG(newmsg)
			}
		}
	}
}
