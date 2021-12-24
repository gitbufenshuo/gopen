package resource

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
