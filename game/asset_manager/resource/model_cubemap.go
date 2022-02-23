package resource

import "github.com/gitbufenshuo/gopen/matmath"

var modelCubemapJSON = `{"Vertices":[
	-0.5,-0.5,0.5
	,0.5,-0.5,0.5
	,0.5,0.5,0.5
	,-0.5,0.5,0.5
	,0.5,-0.5,0.5
	,0.5,-0.5,-0.5
	,0.5,0.5,-0.5
	,0.5,0.5,0.5
	,0.5,-0.5,-0.5
	,-0.5,-0.5,-0.5
	,-0.5,0.5,-0.5
	,0.5,0.5,-0.5
	,-0.5,-0.5,-0.5
	,-0.5,-0.5,0.5
	,-0.5,0.5,0.5
	,-0.5,0.5,-0.5
	,-0.5,0.5,0.5
	,0.5,0.5,0.5
	,0.5,0.5,-0.5
	,-0.5,0.5,-0.5
	,-0.5,-0.5,0.5
	,-0.5,-0.5,-0.5
	,0.5,-0.5,-0.5
	,0.5,-0.5,0.5],"Indices":[0,1,2,0,2,3,4,5,6,4,6,7,8,9,10,8,10,11,12,13,14,12,14,15,16,17,18,16,18,19,20,21,22,20,22,23],"Stripes":[3]}`

func NewCubemapModel() *Model {
	res := NewModel()
	res.ReadFromContent(modelCubemapJSON)
	return res
}

func NewCubemapModel_BySpec(pivot, size matmath.Vec4) *Model {
	rawblock := NewCubemapModel()
	//
	width, height, thick := size.GetValue3()
	xoffset, yoffset, zoffset := pivot.GetValue3()
	var dataCount = 3
	for idx := 0; idx != 24; idx++ {
		if rawblock.Vertices[idx*dataCount+0] < 0 {
			rawblock.Vertices[idx*dataCount+0] = -width / 2
		} else {
			rawblock.Vertices[idx*dataCount+0] = width / 2
		}
		rawblock.Vertices[idx*dataCount+0] += xoffset * (-width / 2)
		//
		if rawblock.Vertices[idx*dataCount+1] < 0 {
			rawblock.Vertices[idx*dataCount+1] = -height / 2
		} else {
			rawblock.Vertices[idx*dataCount+1] = height / 2
		}
		rawblock.Vertices[idx*dataCount+1] += yoffset * (-height / 2)
		//
		if rawblock.Vertices[idx*dataCount+2] < 0 {
			rawblock.Vertices[idx*dataCount+2] = -thick / 2
		} else {
			rawblock.Vertices[idx*dataCount+2] = thick / 2
		}
		rawblock.Vertices[idx*dataCount+2] += zoffset * (-thick / 2)
	}
	return rawblock
}
