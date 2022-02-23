package logic_rotate

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
)

type LogicRotate struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	gb game.GameObjectI
	//
	nowForward matmath.Vec4
}

func NewLogicRotate(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicRotate)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	res.nowForward = matmath.CreateVec4FromStr("0,0,0,0")
	return res
}

func (lj *LogicRotate) Start(gb game.GameObjectI) {
	lj.gb = gb
}

func (lj *LogicRotate) Update(gb game.GameObjectI) {
	zhi := float32(lj.gi.CurFrame) / 10
	//
	x, z := help.Sin(zhi), help.Cos(zhi)
	lj.nowForward.SetIndexValue(0, x)
	lj.nowForward.SetIndexValue(1, x+z)
	lj.nowForward.SetIndexValue(2, z)
	gb.GetTransform().SetForward(lj.nowForward, 1)
}
