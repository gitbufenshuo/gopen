package game

type UILayoutTable struct {
	gi        *GlobalInfo
	id        int
	transform *Transform
	UISpec    UISpec
	//
	ElementWidth  float32 // 相邻元素的间隔宽度
	ElementHeight float32 // 相邻元素的间隔高度
	Rows          int     // 一行不得超过
	Cols          int     // 一列不得超过
	//
	Elements []UICanBeLayout //
}

func NewUILayoutTable(gi *GlobalInfo) *UILayoutTable {
	uilt := new(UILayoutTable)
	uilt.gi = gi
	uilt.transform = NewTransform()
	return uilt
}

func (uilt *UILayoutTable) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return uilt.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	uilt.id = _id[0]
	return uilt.id
}

func (uilt *UILayoutTable) Start() {
}

func (uilt *UILayoutTable) Update() {
	{
		// 根据 UISpec 得到真正要渲染的参数
		widthDeform := uilt.gi.GetWindowWidth() / uilt.gi.UICanvas.DesignWidth
		heightDeform := uilt.gi.GetWindowHeight() / uilt.gi.UICanvas.DesignHeight
		// 1. pos
		{
			posx, posy := uilt.UISpec.LocalPos.GetValue2()
			posrx, posry := uilt.UISpec.PosRelativity.GetValue2()
			// 根据真实分辨率，计算新的位置
			posxNew := posx * widthDeform
			posyNew := posy * heightDeform
			posxNew = (1-posrx)*posx + posrx*posxNew
			posyNew = (1-posry)*posy + posry*posyNew
			uilt.transform.Postion.SetIndexValue(0, posxNew)
			uilt.transform.Postion.SetIndexValue(1, posyNew)
		}
		// 2. scale
		{
			scalex := 1 + uilt.UISpec.SizeRelativity.GetIndexValue(0)*(widthDeform-1)
			scaley := 1 + uilt.UISpec.SizeRelativity.GetIndexValue(1)*(heightDeform-1)
			uilt.transform.Scale.SetValue2(
				scalex,
				scaley,
			)
		}
		// 3. rotate
		{
			// uibutton.transform.Rotation.SetZ(
			// 	float32(uibutton.gi.CurFrame) / 2,
			// )
		}
	}

}

func (uilt *UILayoutTable) SetEles(eles []UICanBeLayout) {
	for _, onebutton := range eles {
		onebutton.GetTransform().SetParent(uilt.transform)
	}
	uilt.Elements = eles
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
	// offx, offy := uilt.UISpec.LocalPos.GetValue2()
	offx, offy := float32(0), float32(0)
	tx, ty := float32(0), float32(0)
	for idx, oneele := range uilt.Elements {
		ty = offy - float32(idx/uilt.Rows)*uilt.ElementHeight
		tx = offx + float32(idx%uilt.Rows)*uilt.ElementWidth
		oneele.GetUISpec().LocalPos.SetValue2(tx, ty)
	}
}

// 每列不得超过 uilt.Cols
func (uilt *UILayoutTable) arrangeByCol() {
	// offx, offy := uilt.UISpec.LocalPos.GetValue2()
	offx, offy := float32(0), float32(0)
	tx, ty := float32(0), float32(0)
	for idx, oneele := range uilt.Elements {
		tx = offx + float32(idx/uilt.Cols)*uilt.ElementWidth
		ty = offy - float32(idx%uilt.Cols)*uilt.ElementHeight
		oneele.GetUISpec().LocalPos.SetValue2(tx, ty)
	}
}
