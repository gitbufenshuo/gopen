package logic_walk

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

type LogicWalk struct {
	gi *game.GlobalInfo
	*supports.NilLogic
}

func NewLogicWalk(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicWalk)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}

func (lj *LogicWalk) Update(gb game.GameObjectI) {
}
func (lj *LogicWalk) Clone() game.LogicSupportI {
	return NewLogicWalk(lj.gi)
}
