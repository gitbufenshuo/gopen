package logic_walk

import (
	"fmt"

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
	if lj.gi.CurFrame%55 == 0 {
		fmt.Println("walk walk")
	}
}
