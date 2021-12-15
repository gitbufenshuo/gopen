package matmath

import (
	"fmt"
	"math"
)

type VECX struct {
	data      []float32
	dimension int
}

func (vecx *VECX) InitDimension(di int) {
	if di == 2 {
		vecx.Init2()
	} else if di == 3 {
		vecx.Init3()
	} else {
		vecx.Init4()
	}
}

func (vecx *VECX) Init2() {
	vecx.dimension = 2
	vecx.data = make([]float32, 2, 2)
}

func (vecx *VECX) Init3() {
	vecx.dimension = 3
	vecx.data = make([]float32, 3, 3)
}

func (vecx *VECX) Init4() {
	vecx.dimension = 4
	vecx.data = make([]float32, 4, 4)
}

///////////////////////////
///////////////////////////
// Methods on VEC3

// check equal
func (self *VECX) EqualsTo(other *VECX) bool {
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

// read-only, return the index value
// node: index begins from 0
func (self *VECX) GetIndexValue(index int) float32 {
	return self.data[index]
}

// set index value
// node: index begins from 0
func (self *VECX) SetIndexValue(index int, value float32) {
	self.data[index] = value
}
func (self *VECX) AddIndexValue(index int, value float32) {
	self.data[index] += value
}
func (self *VECX) SetValue1(value0 float32) {
	self.data[0] = value0
}
func (self *VECX) SetValue2(value0, value1 float32) {
	self.data[0] = value0
	self.data[1] = value1
}
func (self *VECX) SetValue3(value0, value1, value2 float32) {
	self.data[0] = value0
	self.data[1] = value1
	self.data[2] = value2
}
func (self *VECX) SetValue4(value0, value1, value2, value3 float32) {
	self.data[0] = value0
	self.data[1] = value1
	self.data[2] = value2
	self.data[3] = value3
}

// set all elements from the target
// node: index begins from 0
func (self *VECX) CopyValue(target *VECX) {
	for index := 0; index != self.dimension; index++ {
		self.data[index] = target.data[index]
	}
}

// set all elements to 0
func (self *VECX) Clear() {
	for i := 0; i != self.dimension; i++ {
		self.SetIndexValue(i, 0)
	}
}

// read-only, pretty print the mat
func (self *VECX) PrettyShow() {
	for i := 0; i != self.dimension; i++ {
		fmt.Printf("| %v |\n", self.GetIndexValue(i))
	}
}

// read-only, return the dimension about this vecx
func (self *VECX) Di() int {
	return self.dimension
}

// check whether self and other have the same dimension
func (self *VECX) checkHomotype(other *VECX) bool {
	return self.dimension == other.dimension
}

// vec math add, the result will be stored int another new-vecx;
// vec3 and vec4 cannot add, so return nil on that condition;
func (self *VECX) Add(other *VECX, op ...bool) VECX {
	if !self.checkHomotype(other) {
		panic("Add")
	}
	var res VECX
	res.InitDimension(other.Di())
	var sub bool
	if len(op) != 0 && op[0] {
		sub = true
	}
	for i := 0; i != self.dimension; i++ {
		if sub {
			res.data[i] = self.data[i] - other.data[i]
		} else {
			res.data[i] = self.data[i] + other.data[i]
		}
	}
	return res
}

// vec math add in place, will store the result in self
// vec3 and vec4 cannot add, so return with nothing changed on that condition
func (self *VECX) Add_InPlace(other *VECX, op ...bool) {
	if !self.checkHomotype(other) {
		return
	}
	var sub bool
	if len(op) != 0 && op[0] {
		sub = true
	}

	for i := 0; i != self.dimension; i++ {
		if sub {
			self.data[i] = self.data[i] - other.data[i]
		} else {
			self.data[i] = self.data[i] + other.data[i]
		}
	}
}

func (self *VECX) Add2_InPlace(other1, other2 *VECX, op ...bool) {
	if !self.checkHomotype(other1) {
		return
	}
	if !self.checkHomotype(other2) {
		return
	}
	var sub bool
	if len(op) != 0 && op[0] {
		sub = true
	}

	for i := 0; i != self.dimension; i++ {
		if sub {
			self.data[i] = other1.data[i] - other2.data[i]
		} else {
			self.data[i] = other1.data[i] + other2.data[i]
		}
	}
}

func (self *VECX) Sub_InPlace(other *VECX) {
	self.Add_InPlace(other, true)
	return
}

// vec math scale in place, will store the result in self
func (self *VECX) Scale_InPlace(factor float32) {
	for i := 0; i != self.dimension; i++ {
		self.data[i] = self.data[i] * factor
	}
}

// two vec interpolation
// t should be [0,1], but any other real number is supported
func (self *VECX) Interpolation(other *VECX, t float32) VECX {
	if !self.checkHomotype(other) {
		panic("Interpolation")
	}
	var res VECX
	res.InitDimension(other.Di())
	for i := 0; i != self.dimension; i++ {
		res.data[i] = (1-t)*self.data[i] + t*other.data[i]
	}
	return res
}

// grab one col from matx to store in self
func (self *VECX) GrabColToVec(other *MATX, col int) {
	if !other.checkHomotype_vec(self) {
		return
	}
	if col < 1 || col > self.dimension {
		return
	}
	for i := 0; i != self.dimension; i++ {
		self.SetIndexValue(i, other.GetEleByRowAndCol(i+1, col))
	}
}

// mat mul vec, the result will be stored in self;
// mat3 and vec4 cannot mul, so return nil on that condition;
func (self *VECX) RightMul_InPlace(other *MATX) {
	if !other.checkHomotype_vec(self) {
		return
	}
	resarr := make([]float32, self.Di())
	for rowMat := 0; rowMat != other.Di(); rowMat++ {
		for vecidx := 0; vecidx != self.Di(); vecidx++ {
			resarr[rowMat] += self.GetIndexValue(vecidx) * other.GetEleByRowAndCol(rowMat+1, vecidx+1)
		}
	}
	for idx := range resarr {
		self.SetIndexValue(idx, resarr[idx])
	}
	return
}

// self * other, the so-called vec dot
func (self *VECX) Dot(other *VECX) float32 {
	if !self.checkHomotype(other) {
		panic("two vec in-homotype")
	}
	var sum float32
	for i := 0; i < self.dimension; i++ {
		sum += self.data[i] * other.data[i]
	}
	return sum
}

func (self *VECX) Normalize() {
	weight := math.Sqrt(float64(self.data[0]*self.data[0] + self.data[1]*self.data[1] + self.data[2]*self.data[2]))
	if math.Abs(weight) < 0.000001 {
		return
	}
	for idx := range self.data {
		self.data[idx] = self.data[idx] / float32(weight)
	}
}
func (self *VECX) Clone() VECX {
	var res VECX
	res.InitDimension(self.Di())
	//
	res.CopyValue(self)
	return res
}

// left X right
func Vec3Cross(left, right *VECX) VECX {
	var res VECX
	res.Init3()
	(res.data)[0] = (left.data)[1]*(right.data)[2] - (left.data)[2]*(right.data)[1]
	(res.data)[1] = (left.data)[2]*(right.data)[0] - (left.data)[0]*(right.data)[2]
	(res.data)[2] = (left.data)[0]*(right.data)[1] - (left.data)[1]*(right.data)[0]
	return res
}

// generate mat4
// target : the target point coord
func LookAtFrom4(point, target, up *VECX) MATX {
	var left MATX
	left.Init4()
	var right MATX
	right.Init4()
	left.ToIdentity()
	right.ToIdentity()

	// first lets calculate the camera-z and camera-x and camera-y
	// camera-z
	camera_z := point.Add(target, true)
	camera_z.Normalize()
	// camera-x
	camera_x := Vec3Cross(up, &camera_z)
	camera_x.Normalize()

	// camera-y
	camera_y := Vec3Cross(&camera_z, &camera_x)
	camera_y.Normalize()
	// deal with the left mat4
	left.SetEleByRowAndCol(1, 1, (camera_x.data)[0])
	left.SetEleByRowAndCol(1, 2, (camera_x.data)[1])
	left.SetEleByRowAndCol(1, 3, (camera_x.data)[2])

	left.SetEleByRowAndCol(2, 1, (camera_y.data)[0])
	left.SetEleByRowAndCol(2, 2, (camera_y.data)[1])
	left.SetEleByRowAndCol(2, 3, (camera_y.data)[2])

	left.SetEleByRowAndCol(3, 1, (camera_z.data)[0])
	left.SetEleByRowAndCol(3, 2, (camera_z.data)[1])
	left.SetEleByRowAndCol(3, 3, (camera_z.data)[2])

	// deal with the right mat4
	right.SetEleByRowAndCol(1, 4, -(point.data)[0])
	right.SetEleByRowAndCol(2, 4, -(point.data)[1])
	right.SetEleByRowAndCol(3, 4, -(point.data)[2])

	// left * right --> view mat4
	right.RightMul_InPlace(&left)
	return right
}
