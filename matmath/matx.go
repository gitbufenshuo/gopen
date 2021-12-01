package matmath

import (
	"fmt"
	"math"
	"sync"
)

var global_mat_pool map[int]*sync.Pool

func init() {
	global_mat_pool = make(map[int]*sync.Pool)
	for i := 2; i != 5; i++ {
		var newfunc = func(di int) func() interface{} {
			switch di {
			case 2:
				return newMAT2
			case 3:
				return newMAT3
			case 4:
				return newMAT4
			default:
				panic("MATX only supports [2, 3, 4] dimension")
			}
		}(i)
		global_sync_pool := &sync.Pool{
			New: newfunc,
		}
		global_mat_pool[i] = global_sync_pool
	}

}

// element in the data is arranged in column.
// For example: int a mat2 , data[0] == mat2[1][1], data[1] == mat2[2][1]
type MATX struct {
	data      []float32
	dimension int
}

func newMAT2() interface{} {
	var matx MATX
	matx.dimension = 2
	matx.data = make([]float32, 4, 4)
	return &matx
}
func newMAT3() interface{} {
	var matx MATX
	matx.dimension = 3
	matx.data = make([]float32, 9, 9)
	return &matx
}
func newMAT4() interface{} {
	var matx MATX
	matx.dimension = 4
	matx.data = make([]float32, 16, 16)
	return &matx
}

// this is the only function that you should new-a-MATX
// the data inside the result is dirty, don't think they are zero-initialized
func GetMATX(di int) *MATX {
	return global_mat_pool[di].Get().(*MATX)
}

// when you dont need one MATX anymore, you should call this function
func DontNeedMATXAnyMore(matx *MATX) {
	if matx == nil {
		return
	}
	global_mat_pool[matx.dimension].Put(matx)
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
func (self *MATX) Add(other *MATX) *MATX {
	if !self.checkHomotype(other) {
		return nil
	}
	res := GetMATX(self.dimension)
	for i := 0; i != self.dimension; i++ {
		res.data[i] = self.data[i] + other.data[i]
	}
	return res
}

// mat math add in place, the result will be stored int self;
// mat3 and mat4 cannot add, so return with nothing changed on that condition
func (self *MATX) Add_InPlace(other *MATX) {
	if !self.checkHomotype(other) {
		return
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
	temp_vec := GetVECX(self.dimension)
	for i := 0; i != self.dimension; i++ {
		temp_vec.GrabColToVec(self, i+1)
		temp_vec.RightMul_InPlace(other)
		self.SetByCol_InPlace(temp_vec, i+1)
	}
	DontNeedVECXAnyMore(temp_vec)
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

// in-place
func (self *MATX) Rotate4(rotation *VECX) {
	if self.dimension != 4 {
		panic("Rotate4 dimension doesn't match")
	}
	res := GetMATX(4)
	res.ToIdentity()
	defer DontNeedMATXAnyMore(res)
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

	(res.data)[0] = cosy * cosz
	(res.data)[4] = -(cosy * sinz)
	(res.data)[8] = siny

	(res.data)[1] = cosx*sinz + sinx*siny*cosz
	(res.data)[5] = cosx*cosz - sinx*siny*sinz
	(res.data)[9] = -(sinx * cosy)

	(res.data)[2] = sinx*sinz - cosx*siny*cosz
	(res.data)[6] = sinx*cosz + cosx*siny*sinz
	(res.data)[10] = cosx * cosy

	self.RightMul_InPlace(res)
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
	res := GetMATX(4)
	res.ToIdentity()
	defer DontNeedMATXAnyMore(res)

	(res.data)[12] = translate.GetIndexValue(0)
	(res.data)[13] = translate.GetIndexValue(1)
	(res.data)[14] = translate.GetIndexValue(2)
	self.RightMul_InPlace(res)
}

func Homoz4(z float32) *MATX {
	res := GetMATX(4)
	res.ToIdentity()
	p := -(1 / z)
	res.SetEleByRowAndCol(3, 3, 0.0)
	res.SetEleByRowAndCol(4, 4, 0.0)
	res.SetEleByRowAndCol(4, 3, p)
	res.SetEleByRowAndCol(3, 4, p)
	return res
}

// @title    Perspective
// @description   获取一个 Perspective 的矩阵
// @auth      onebook
// @param     near near plane
// @param     far far plane
// @param     fov
func Perspective(near, far, fov float32) *MATX {
	res := GetMATX(4)
	res.ToIdentity()
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
