package logic_jump

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

type LogicJump struct {
	gi *game.GlobalInfo
	*supports.NilLogic
}

func NewLogicJump(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicJump)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}

func (lj *LogicJump) Update(gb game.GameObjectI) {

}
func (lj *LogicJump) Clone() game.LogicSupportI {
	return NewLogicJump(lj.gi)
}
