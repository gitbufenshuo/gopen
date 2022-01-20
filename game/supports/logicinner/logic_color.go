package logicinner

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

type LogicColorControl struct {
	*supports.NilLogic
	Color [3]float32
}

func NewLogicColorControl() *LogicColorControl {
	res := new(LogicColorControl)
	res.NilLogic = supports.NewNilLogic()
	res.Color[0], res.Color[1], res.Color[2] = 1, 1, 1
	return res
}

func (lcc *LogicColorControl) OnDraw(gb game.GameObjectI) {
	rs := gb.GetRenderSupport()
	rs.ShaderOP().SetUniform3f("u_Color", lcc.Color[0], lcc.Color[1], lcc.Color[2])
}
func (lcc *LogicColorControl) OnDrawFinish(gb game.GameObjectI) {
	rs := gb.GetRenderSupport()
	rs.ShaderOP().SetUniform3f("u_Color", 0, 0, 0)
}
func (lcc *LogicColorControl) Clone() game.LogicSupportI {
	return NewLogicColorControl()
}
