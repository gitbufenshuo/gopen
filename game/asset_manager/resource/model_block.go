package resource

import "github.com/gitbufenshuo/gopen/matmath"

var modelBlockJSON = `{"Vertices":[-0.5,-0.5,0.5,0,0,0,0,1,0.5,-0.5,0.5,1,0,0,0,1,0.5,0.5,0.5,1,1,0,0,1,-0.5,0.5,0.5,0,1,0,0,1,0.5,-0.5,0.5,0,0,1,0,0,0.5,-0.5,-0.5,1,0,1,0,0,0.5,0.5,-0.5,1,1,1,0,0,0.5,0.5,0.5,0,1,1,0,0,0.5,-0.5,-0.5,0,0,0,0,-1,-0.5,-0.5,-0.5,1,0,0,0,-1,-0.5,0.5,-0.5,1,1,0,0,-1,0.5,0.5,-0.5,0,1,0,0,-1,-0.5,-0.5,-0.5,0,0,-1,0,0,-0.5,-0.5,0.5,1,0,-1,0,0,-0.5,0.5,0.5,1,1,-1,0,0,-0.5,0.5,-0.5,0,1,-1,0,0,-0.5,0.5,0.5,0,0,0,1,0,0.5,0.5,0.5,1,0,0,1,0,0.5,0.5,-0.5,1,1,0,1,0,-0.5,0.5,-0.5,0,1,0,1,0,-0.5,-0.5,0.5,0,0,0,-1,0,-0.5,-0.5,-0.5,1,0,0,-1,0,0.5,-0.5,-0.5,1,1,0,-1,0,0.5,-0.5,0.5,0,1,0,-1,0],"Indices":[0,1,2,0,2,3,4,5,6,4,6,7,8,9,10,8,10,11,12,13,14,12,14,15,16,17,18,16,18,19,20,21,22,20,22,23],"Stripes":[3,2,3]}`
var BlockModel *Model
var PlaneModel *Model

func init() {
	BlockModel = NewModel()
	BlockModel.ReadFromContent(modelBlockJSON)
	//
	PlaneModel = NewModel()
	PlaneModel.ReadFromContent(modelBlockJSON)
	for idx := range PlaneModel.Vertices {
		if idx%8 == 0 || idx%8 == 2 {
			PlaneModel.Vertices[idx] *= 10
		}
	}
}

func NewBlockModel() *Model {
	res := NewModel()
	res.ReadFromContent(modelBlockJSON)
	return res
}

func NewBlockModel_BySpec(pivot, size matmath.Vec4) *Model {
	rawblock := NewBlockModel()
	//
	width, height, thick := size.GetValue3()
	xoffset, yoffset, zoffset := pivot.GetValue3()
	for idx := 0; idx != 24; idx++ {
		if rawblock.Vertices[idx*8+0] < 0 {
			rawblock.Vertices[idx*8+0] = -width / 2
		} else {
			rawblock.Vertices[idx*8+0] = width / 2
		}
		rawblock.Vertices[idx*8+0] += xoffset
		//
		if rawblock.Vertices[idx*8+1] < 0 {
			rawblock.Vertices[idx*8+1] = -height / 2
		} else {
			rawblock.Vertices[idx*8+1] = height / 2
		}
		rawblock.Vertices[idx*8+1] += yoffset
		//
		if rawblock.Vertices[idx*8+2] < 0 {
			rawblock.Vertices[idx*8+2] = -thick / 2
		} else {
			rawblock.Vertices[idx*8+2] = thick / 2
		}
		rawblock.Vertices[idx*8+2] += zoffset
	}
	return rawblock
}
