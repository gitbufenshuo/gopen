package resource

var modelQuadJSON = `{"Vertices":[
-0.1,-0.1,0,0,0,
0.1,-0.1,0,1,0,
0.1,0.1,0,1,1,
-0.1,0.1,0,0,1
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
