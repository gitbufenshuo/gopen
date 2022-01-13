package logicinner

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
	"github.com/gitbufenshuo/gopen/matmath"
)

type LogicRotate struct {
	*supports.NilLogic
	RotateValue *matmath.Vec4
}

func NewLogicRotate(data string) *LogicRotate {
	res := new(LogicRotate)
	//
	res.NilLogic = supports.NewNilLogic()
	a := matmath.CreateVec4FromStr(data)
	res.RotateValue = &a
	return res
}

func (lr *LogicRotate) Update(gb game.GameObjectI) {
	x, y, z := lr.RotateValue.GetValue3()
	gb.GetTransform().Rotation.AddIndexValue(0, x)
	gb.GetTransform().Rotation.AddIndexValue(1, y)
	gb.GetTransform().Rotation.AddIndexValue(2, z)
}
