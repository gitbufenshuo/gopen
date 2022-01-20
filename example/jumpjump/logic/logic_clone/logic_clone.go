package logic_clone

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

type LogicClone struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	NoClone bool // 不允许clone吗
}

func newLogicClone(gi *game.GlobalInfo) *LogicClone {
	res := new(LogicClone)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	return res
}

func NewLogicClone(gi *game.GlobalInfo) game.LogicSupportI {
	return newLogicClone(gi)
}

func (lc *LogicClone) Update(gb game.GameObjectI) {
	if lc.gi.CurFrame == 100 && !lc.NoClone {
		fmt.Println("LogicClone", gb.ID_sg())
		newgb := lc.gi.InstantiateGameObject(gb)
		newtr := newgb.GetTransform()
		newtr.Postion.SetValue2(-3, 2)
	}
}
func (lc *LogicClone) Clone() game.LogicSupportI {
	newlc := newLogicClone(lc.gi)
	newlc.NoClone = true
	return newlc
}
