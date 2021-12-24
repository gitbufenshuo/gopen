package matmath

import "fmt"

type Vec4 struct {
	data [4]float32
}

func CreateVec4(x, y, z, w float32) Vec4 {
	return Vec4{
		data: [4]float32{x, y, z, w},
	}
}
func (vec4 *Vec4) PrettyShow() {
	fmt.Printf("(")
	for i := 0; i != 4; i++ {
		fmt.Printf(" %v ", vec4.data[i])
	}
	fmt.Printf(")\n")
}

func (vec4 Vec4) LeftMulMAT(mat MATX) Vec4 {
	x, y, z, w := vec4.data[0], vec4.data[1], vec4.data[2], vec4.data[3]
	//
	return CreateVec4(
		x*mat.GetEleByRowAndCol(1, 1)+y*mat.GetEleByRowAndCol(1, 2)+z*mat.GetEleByRowAndCol(1, 3)+w*mat.GetEleByRowAndCol(1, 4),
		x*mat.GetEleByRowAndCol(2, 1)+y*mat.GetEleByRowAndCol(2, 2)+z*mat.GetEleByRowAndCol(2, 3)+w*mat.GetEleByRowAndCol(2, 4),
		x*mat.GetEleByRowAndCol(3, 1)+y*mat.GetEleByRowAndCol(3, 2)+z*mat.GetEleByRowAndCol(3, 3)+w*mat.GetEleByRowAndCol(3, 4),
		x*mat.GetEleByRowAndCol(4, 1)+y*mat.GetEleByRowAndCol(4, 2)+z*mat.GetEleByRowAndCol(4, 3)+w*mat.GetEleByRowAndCol(4, 4),
	)
}
