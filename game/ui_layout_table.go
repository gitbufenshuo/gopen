package game

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game/common"
)

type UILayoutTable struct {
	gi        *GlobalInfo
	transform *common.Transform
	UISpec    UISpec
	//
	ElementWidth  float32 // 相邻元素的间隔宽度
	ElementHeight float32 // 相邻元素的间隔高度
	Rows          int     // 一行不得超过
	Cols          int     // 一列不得超过
	//
	Elements []*UIButton //
}

func NewUILayoutTable(gi *GlobalInfo) *UILayoutTable {
	uilt := new(UILayoutTable)
	uilt.gi = gi
	uilt.transform = common.NewTransform()
	return uilt
}

// 排列一次，根据所有的参数，算出最终的各个元素的参数
// 主要是改变元素的 UISpec.LocalPos
func (uilt *UILayoutTable) Arrange() {
	if uilt.Rows+uilt.Cols == 0 {
		return
	}
	if uilt.Rows > 0 {
		uilt.arrangeByRow()
		return
	}
	if uilt.Cols > 0 {
		uilt.arrangeByCol()
		return
	}
}

// 每行不得超过 uilt.Rows
func (uilt *UILayoutTable) arrangeByRow() {
	offx, offy := uilt.UISpec.LocalPos.GetValue2()
	tx, ty := float32(0), float32(0)
	for idx, oneele := range uilt.Elements {
		ty = offy - float32(idx/uilt.Rows)*uilt.ElementHeight
		tx = offx + float32(idx%uilt.Rows)*uilt.ElementWidth
		oneele.UISpec.LocalPos.SetValue2(tx, ty)
		fmt.Println(offx, offy)
	}
}

// 每列不得超过 uilt.Cols
func (uilt *UILayoutTable) arrangeByCol() {
	offx, offy := uilt.UISpec.LocalPos.GetValue2()
	tx, ty := float32(0), float32(0)
	for idx, oneele := range uilt.Elements {
		tx = offx + float32(idx/uilt.Cols)*uilt.ElementWidth
		ty = offy - float32(idx%uilt.Cols)*uilt.ElementHeight
		oneele.UISpec.LocalPos.SetValue2(tx, ty)
	}
}
