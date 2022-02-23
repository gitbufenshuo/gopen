package resource

import (
	"github.com/gitbufenshuo/gopen/matmath"
)

var modelCubemapJSON = `{"Vertices":[
    -1.0,  1.0, -1.0,
    -1.0, -1.0, -1.0,
     1.0, -1.0, -1.0,
     1.0, -1.0, -1.0,
     1.0,  1.0, -1.0,
    -1.0,  1.0, -1.0,

    -1.0, -1.0,  1.0,
    -1.0, -1.0, -1.0,
    -1.0,  1.0, -1.0,
    -1.0,  1.0, -1.0,
    -1.0,  1.0,  1.0,
    -1.0, -1.0,  1.0,

     1.0, -1.0, -1.0,
     1.0, -1.0,  1.0,
     1.0,  1.0,  1.0,
     1.0,  1.0,  1.0,
     1.0,  1.0, -1.0,
     1.0, -1.0, -1.0,

    -1.0, -1.0,  1.0,
    -1.0,  1.0,  1.0,
     1.0,  1.0,  1.0,
     1.0,  1.0,  1.0,
     1.0, -1.0,  1.0,
    -1.0, -1.0,  1.0,

    -1.0,  1.0, -1.0,
     1.0,  1.0, -1.0,
     1.0,  1.0,  1.0,
     1.0,  1.0,  1.0,
    -1.0,  1.0,  1.0,
    -1.0,  1.0, -1.0,

    -1.0, -1.0, -1.0,
    -1.0, -1.0,  1.0,
     1.0, -1.0, -1.0,
     1.0, -1.0, -1.0,
    -1.0, -1.0,  1.0,
     1.0, -1.0,  1.0
],"Indices":[],"Stripes":[3]}`

func NewCubemapModel() *Model {
	res := NewModel()
	res.ReadFromContent(modelCubemapJSON)
	return res
}

func NewCubemapModel_BySpec(pivot, size matmath.Vec4) *Model {
	rawblock := NewCubemapModel()
	rawblock.Indices = make([]uint32, 36)
	for idx := range rawblock.Indices {
		rawblock.Indices[idx] = uint32(idx)
	}
	return rawblock
	// rawblock.Indices = []uint32{0, 1, 2, 0, 2, 3, 10, 9, 8, 11, 10, 8}
	// return rawblock
	//
	width, height, thick := size.GetValue3()
	xoffset, yoffset, zoffset := pivot.GetValue3()
	var dataCount = 3
	for idx := 0; idx != 24; idx++ {
		// rawblock.Vertices[idx*dataCount+2] *= -1
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
