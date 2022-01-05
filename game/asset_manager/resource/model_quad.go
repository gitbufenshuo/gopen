package resource

import "github.com/gitbufenshuo/gopen/matmath"

var modelQuadJSON = `{"Vertices":[
-1,-1,0,0,0,
1,-1,0,1,0,
1,1,0,1,1,
-1,1,0,0,1
],"Indices":[0,1,2,0,2,3],"Stripes":[3,2]}`
var QuadModel *Model

func init() {
	QuadModel = NewModel()
	QuadModel.ReadFromContent(modelQuadJSON)
}

func NewQuadModel() *Model {
	res := NewModel()
	res.ReadFromContent(modelQuadJSON)
	return res
}

func NewQuadModel_LeftALign() *Model {
	res := NewModel()
	res.ReadFromContent(modelQuadJSON)
	for idx := 0; idx != 4; idx++ {
		if res.Vertices[idx*5+0] < 0 {
			res.Vertices[idx*5+0] = 0
		} else {
			res.Vertices[idx*5+0] *= 2
		}
	}
	return res
}

func NewQuadModel_BySpec(pivot matmath.Vec4, width, height float32) *Model {
	res := NewModel()
	res.ReadFromContent(modelQuadJSON)
	xoffset := -(width / 2) * pivot.GetIndexValue(0)
	yoffset := -(height / 2) * pivot.GetIndexValue(1)
	for idx := 0; idx != 4; idx++ {
		if res.Vertices[idx*5+0] < 0 {
			res.Vertices[idx*5+0] = -width / 2
		} else {
			res.Vertices[idx*5+0] = width / 2
		}
		res.Vertices[idx*5+0] += xoffset
		if res.Vertices[idx*5+1] < 0 {
			res.Vertices[idx*5+1] = -height / 2
		} else {
			res.Vertices[idx*5+1] = height / 2
		}
		res.Vertices[idx*5+1] += yoffset
	}

	return res
}
