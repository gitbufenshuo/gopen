package matmath

import (
	"fmt"
	"math"
)

// element in the data is arranged in column.
// For example: int a mat2 , data[0] == mat2[1][1], data[1] == mat2[2][1]
type MATX struct {
	data      []float32
	dimension int
}

func (matx *MATX) InitDimension(di int) {
	if di == 2 {
		matx.Init2()
	} else if di == 3 {
		matx.Init3()
	} else {
		matx.Init4()
	}
}

func (matx *MATX) Init2() {
	matx.dimension = 2
	matx.data = make([]float32, 4, 4)
}

func (matx *MATX) Init3() {
	matx.dimension = 3
	matx.data = make([]float32, 9, 9)
}

func (matx *MATX) Init4() {
	matx.dimension = 4
	matx.data = make([]float32, 16, 16)
}

// read-only, pretty print the mat
func (self *MATX) PrettyShow(logs ...string) {
	fmt.Println(">>><<<", logs)
	for row := 1; row <= self.dimension; row++ {
		fmt.Printf("|")
		for col := 1; col <= self.dimension; col++ {
			fmt.Printf(" %v ", self.GetEleByRowAndCol(row, col))
		}
		fmt.Printf("|\n")
	}
}

// check equal
func (self *MATX) EqualsTo(other *MATX) bool {
	if self == other {
		return true
	}
	if !self.checkHomotype(other) {
		return false
	}
	for idx := range self.data {
		if !(math.Abs(float64(self.data[idx]-other.data[idx])) < 0.00000001) {
			return false
		}
	}
	return true
}

// read-only, return the dimension about this MATX
func (self *MATX) Di() int {
	return self.dimension
}

// check whether self and other have the same dimension
func (self *MATX) checkHomotype(other *MATX) bool {
	return self.dimension == other.dimension
}

// check whether self and other vecx have the same dimension
func (self *MATX) checkHomotype_vec(other *VECX) bool {
	return self.dimension == other.dimension
}

// mat math add, the result will be stored int another new-matx;
// mat3 and mat4 cannot add, so return nil on that condition;
func (self *MATX) Add(other *MATX) MATX {
	if !self.checkHomotype(other) {
		panic("mat add dimension")
	}
	var res MATX
	p := &res
	p.InitDimension(other.Di())
	p = nil
	for i := 0; i != self.dimension; i++ {
		res.data[i] = self.data[i] + other.data[i]
	}
	return res
}

// mat math add in place, the result will be stored int self;
// mat3 and mat4 cannot add, so return with nothing changed on that condition
func (self *MATX) Add_InPlace(other *MATX) {
	if !self.checkHomotype(other) {
		panic("MATX.Add_InPlace")
	}
	for i := 0; i != self.dimension; i++ {
		self.data[i] = self.data[i] + other.data[i]
	}
	return
}

// matx[row][ele] = value
func (self *MATX) SetEleByRowAndCol(row, col int, value float32) {
	if row > self.dimension || row < 1 {
		return
	}
	if row > self.dimension || row < 1 {
		return
	}
	eleIndex := (col-1)*self.dimension + row - 1
	self.data[eleIndex] = value
}

// return matx[row][ele]
func (self *MATX) GetEleByRowAndCol(row, col int) float32 {
	if row > self.dimension || row < 1 {
		panic("out of index")
	}
	if col > self.dimension || col < 1 {
		panic("out of index")
	}
	eleIndex := (col-1)*self.dimension + row - 1
	return self.data[eleIndex]
}

// set one col with a vec in place
// note: col begins from 1 for humanity
func (self *MATX) SetByCol_InPlace(other *VECX, col int) {
	if !self.checkHomotype_vec(other) {
		return
	}
	if col < 1 || col > self.dimension {
		return
	}
	for i := 0; i != self.dimension; i++ {
		self.SetEleByRowAndCol(i+1, col, other.GetIndexValue(i))
	}
}

// self = other * self
// mat3 and vec4 cannot mul, so return nil on that condition;
func (self *MATX) RightMul_InPlace(other *MATX) {
	if !self.checkHomotype(other) {
		return
	}
	var temp_vec VECX
	temp_vec.InitDimension(self.Di())
	for i := 0; i != self.dimension; i++ {
		temp_vec.GrabColToVec(self, i+1)
		temp_vec.RightMul_InPlace(other)
		self.SetByCol_InPlace(&temp_vec, i+1)
	}
	return
}

// in-place
func (self *MATX) ToIdentity() {
	modula := self.dimension + 1
	for idx := range self.data {
		if idx%modula == 0 {
			self.data[idx] = 1
		} else {
			self.data[idx] = 0
		}
	}
}

// return the first element address
func (self *MATX) Address() *float32 {
	return &(self.data[0])
}

func (self *MATX) Data() []float32 {
	return self.data
}

// in-place
func (self *MATX) Scale4(scale *VECX) {
	if self.dimension != 4 {
		panic("Rotate4 dimension doesn't match")
	}
	self.data[0] = scale.data[0]
	self.data[5] = scale.data[1]
	self.data[10] = scale.data[1]
}

// in-place
func (self *MATX) Rotate4(rotation *VECX) {
	if self.dimension != 4 {
		panic("Rotate4 dimension doesn't match")
	}
	var helperMat MATX
	helperMat.Init4()
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

	self.RightMul_InPlace(&helperMat)
}

// in-place
func (self *MATX) Translate4(translate *VECX) {
	// ST_MAT4 *res = NewMat4(0);
	// Mat4SetValue(res, 1, 4, x_value);
	// Mat4SetValue(res, 2, 4, y_value);
	// Mat4SetValue(res, 3, 4, z_value);
	// ST_MAT4 *shouldReturn = MatMat4(res, mat4);
	// Mat4Free(res);
	// return shouldReturn;
	///////////////////////
	if self.dimension != 4 {
		panic("Translate4 dimension doesn't match")
	}
	var helperMat MATX
	p := &helperMat
	p.Init4()
	p.ToIdentity()
	(p.data)[12] = translate.GetIndexValue(0)
	(p.data)[13] = translate.GetIndexValue(1)
	(p.data)[14] = translate.GetIndexValue(2)
	self.RightMul_InPlace(p)
}

// @title    Perspective
// @description   获取一个 Perspective 的矩阵
// @auth      onebook
// @param     near near plane
// @param     far far plane
// @param     fov
func Perspective(near, far, fov float32) MATX {
	var res MATX
	p := &res
	p.Init4()
	p.ToIdentity()
	//
	topdown := float32(math.Tan(float64(fov/2))) * near
	leftright := topdown // cause aspect is always 1
	res.SetEleByRowAndCol(1, 1, 2*near/(leftright))
	res.SetEleByRowAndCol(2, 2, 2*near/(leftright))
	res.SetEleByRowAndCol(3, 3, -(far+near)/(far-near))
	res.SetEleByRowAndCol(3, 4, 2*far*near/(near-far))
	res.SetEleByRowAndCol(4, 3, 2*far*near/(near-far))
	return res
}
