package matmath

import (
	"fmt"
	"math"
)

// element in the data is arranged in column.
// For example: int a mat2 , data[0] == mat2[1][1], data[1] == mat2[2][1]
type MAT4 struct {
	data [16]float32
}

// read-only, pretty print the mat
func (mat4 *MAT4) PrettyShow(logs ...string) {
	fmt.Println(">>><<<", logs)
	for row := 1; row <= 4; row++ {
		fmt.Printf("|")
		for col := 1; col <= 4; col++ {
			fmt.Printf(" %v ", mat4.GetEleByRowAndCol(row, col))
		}
		fmt.Printf("|\n")
	}
}

// check equal
func (mat4 *MAT4) EqualsTo(other *MAT4) bool {
	if mat4 == other {
		return true // 地址是一样
	}

	for idx := range mat4.data {
		if !(math.Abs(float64(mat4.data[idx]-other.data[idx])) < 0.00000001) {
			return false
		}
	}
	return true
}

// mat math add, the result will be stored int another new-matx;
// mat3 and mat4 cannot add, so return nil on that condition;
func (mat4 *MAT4) Add(other *MAT4) MATX {
	var res MATX
	for i := 0; i != 4; i++ {
		res.data[i] = mat4.data[i] + other.data[i]
	}
	return res
}

// mat math add in place, the result will be stored int self;
// mat3 and mat4 cannot add, so return with nothing changed on that condition
func (mat4 *MAT4) Add_InPlace(other *MAT4) {
	for i := 0; i != 4; i++ {
		mat4.data[i] = mat4.data[i] + other.data[i]
	}
	return
}

// matx[row][ele] = value
func (mat4 *MAT4) SetEleByRowAndCol(row, col int, value float32) {
	eleIndex := (col-1)*4 + row - 1
	mat4.data[eleIndex] = value
}

// return matx[row][ele]
func (mat4 *MAT4) GetEleByRowAndCol(row, col int) float32 {
	eleIndex := (col-1)*4 + row - 1
	return mat4.data[eleIndex]
}

// set one col with a vec in place
// note: col begins from 1 for humanity
func (mat4 *MAT4) SetByCol_InPlace(other *Vec4, col int) {
	for i := 0; i != 4; i++ {
		mat4.SetEleByRowAndCol(i+1, col, other.GetIndexValue(i))
	}
}

// self = other * self
// mat3 and vec4 cannot mul, so return nil on that condition;
func (mat4 *MAT4) RightMul_InPlace(other *MAT4) {
	var temp_vec Vec4
	for i := 0; i != 4; i++ {
		temp_vec.GrabColToVec(mat4, i+1)
		temp_vec.RightMul_InPlace(other)
		mat4.SetByCol_InPlace(&temp_vec, i+1)
	}
	return
}

// in-place
func (mat4 *MAT4) ToIdentity() {
	modula := 5
	for idx := range mat4.data {
		if idx%modula == 0 {
			mat4.data[idx] = 1
		} else {
			mat4.data[idx] = 0
		}
	}
}

// return the first element address
func (mat4 *MAT4) Address() *float32 {
	return &(mat4.data[0])
}

func (mat4 *MAT4) Data() [16]float32 {
	return mat4.data
}

// in-place
func (mat4 *MAT4) Scale(scale *Vec4) {
	mat4.data[0] = scale.data[0]
	mat4.data[5] = scale.data[1]
	mat4.data[10] = scale.data[1]
}

// in-place
func (mat4 *MAT4) Rotate(rotation *Vec4) {
	var helperMat MAT4
	helperMat.ToIdentity()
	//
	x_degree := rotation.GetIndexValue(0)
	y_degree := rotation.GetIndexValue(1)
	z_degree := rotation.GetIndexValue(2)
	//
	cosy := float32(math.Cos(float64((3.141592653 * y_degree) / 180.0)))
	cosz := float32(math.Cos(float64((3.141592653 * z_degree) / 180.0)))
	sinz := float32(math.Sin(float64((3.141592653 * z_degree) / 180.0)))
	siny := float32(math.Sin(float64((3.141592653 * y_degree) / 180.0)))
	cosx := float32(math.Cos(float64((3.141592653 * x_degree) / 180.0)))
	sinx := float32(math.Sin(float64((3.141592653 * x_degree) / 180.0)))

	(helperMat.data)[0] = cosy * cosz
	(helperMat.data)[4] = -(cosy * sinz)
	(helperMat.data)[8] = siny

	(helperMat.data)[1] = cosx*sinz + sinx*siny*cosz
	(helperMat.data)[5] = cosx*cosz - sinx*siny*sinz
	(helperMat.data)[9] = -(sinx * cosy)

	(helperMat.data)[2] = sinx*sinz - cosx*siny*cosz
	(helperMat.data)[6] = sinx*cosz + cosx*siny*sinz
	(helperMat.data)[10] = cosx * cosy

	mat4.RightMul_InPlace(&helperMat)
}

// in-place
func (mat4 *MAT4) Translate4(translate *Vec4) {
	var helperMat MAT4
	helperMat.ToIdentity()
	(helperMat.data)[12] = translate.GetIndexValue(0)
	(helperMat.data)[13] = translate.GetIndexValue(1)
	(helperMat.data)[14] = translate.GetIndexValue(2)
	mat4.RightMul_InPlace(&helperMat)
}

// @title    Perspective
// @description   获取一个 Perspective 的矩阵
// @auth      onebook
// @param     near near plane
// @param     far far plane
// @param     fov
func GenPerspectiveMat4(near, far, fov float32, aspect float32) MAT4 {
	var res MAT4
	res.ToIdentity()
	//
	topdown := float32(math.Tan(float64(fov/2))) * near
	leftright := topdown * aspect // cause aspect is always 1
	res.SetEleByRowAndCol(1, 1, 2*near/(leftright))
	res.SetEleByRowAndCol(2, 2, 2*near/(topdown))
	res.SetEleByRowAndCol(3, 3, -(far+near)/(far-near))
	res.SetEleByRowAndCol(3, 4, 2*far*near/(near-far))
	res.SetEleByRowAndCol(4, 3, -1)
	return res
}

func GenOrthographicMat4(near, far float32, topdown, aspect float32) MAT4 {
	var res MAT4
	res.ToIdentity()
	//
	leftright := topdown * aspect
	res.SetEleByRowAndCol(1, 1, 2/(leftright))
	res.SetEleByRowAndCol(2, 2, 2/(topdown))
	res.SetEleByRowAndCol(3, 3, 2/(near-far))
	res.SetEleByRowAndCol(3, 4, -(far+near)/(far-near))
	return res
}
