package logic_follow

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/help"
)

type LogicFollow struct {
	gi *game.GlobalInfo
	*supports.NilLogic
	//
	logicposx, logicposy, logicposz int64
	factor                          float32
}

func NewLogicFollow(gi *game.GlobalInfo) game.LogicSupportI {
	res := new(LogicFollow)
	//
	res.NilLogic = supports.NewNilLogic()
	res.gi = gi
	res.factor = 5
	return res
}
func (lj *LogicFollow) syncLogicPosY(gb game.GameObjectI) {
	nowposx, nowposy, nowposz := gb.GetTransform().Postion.GetValue3()
	nowposx += (float32(lj.logicposx)/1000 - nowposx) / lj.factor
	nowposy += (float32(lj.logicposy)/1000 - nowposy) / lj.factor
	nowposz += (float32(lj.logicposz)/1000 - nowposz) / lj.factor
	gb.GetTransform().Postion.SetValue3(
		nowposx, nowposy, nowposz,
	)
}

func (lj *LogicFollow) Move(x, y, z, chazhi int64) { // 逻辑相关
	lj.logicposx = help.Int64To(lj.logicposx, x, chazhi)
	lj.logicposy = help.Int64To(lj.logicposy, y, chazhi)
	lj.logicposz = help.Int64To(lj.logicposz, z, chazhi)
}

func (lj *LogicFollow) Update(gb game.GameObjectI) {
	lj.syncLogicPosY(gb) // 不影响逻辑
}

func (lj *LogicFollow) Clone() game.LogicSupportI {
	return NewLogicFollow(lj.gi)
}
