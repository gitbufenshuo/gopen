package logic_clone

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

type LogicClone struct {
	gi *game.GlobalInfo
	*supports.NilLogic
}

func NewLogicClone(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicClone)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}

func (lc *LogicClone) Update(gb game.GameObjectI) {
	if lc.gi.CurFrame == 100 {
		fmt.Println("LogicClone", gb.ID_sg())
		newgb := lc.gi.InstantiateGameObject(gb)
		newtr := newgb.GetTransform()
		newtr.Postion.SetValue2(-3, 2)
	}
}
func (lc *LogicClone) Clone() game.LogicSupportI {
	return NewLogicClone(lc.gi)
}
