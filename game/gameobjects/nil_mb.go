package gameobjects

type NilManageObject struct {
	ID int
}

func NewNilManageObject() *NilManageObject {
	res := new(NilManageObject)
	return res
}

func (nmo *NilManageObject) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return nmo.ID
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	nmo.ID = _id[0]
	return nmo.ID
}

func (nmo *NilManageObject) Start() {
}

func (nmo *NilManageObject) Update() {
}
