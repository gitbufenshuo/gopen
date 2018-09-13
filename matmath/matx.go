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
	global_mat_pool[matx.dimension].Put(matx)
}

// read-only, pretty print the mat
func (self *MATX) PrettyShow() {
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
