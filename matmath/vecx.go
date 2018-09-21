package matmath

import (
	"fmt"
	"math"
	"sync"
)

var global_vec_pool map[int]*sync.Pool

func init() {
	global_vec_pool = make(map[int]*sync.Pool)
	for i := 2; i != 5; i++ {
		var newfunc = func(di int) func() interface{} {
			switch di {
			case 2:
				return newVEC2
			case 3:
				return newVEC3
			case 4:
				return newVEC4
			default:
				panic("VECX only supports [2, 3, 4] dimension")
			}
		}(i)
		global_sync_pool := &sync.Pool{
			New: newfunc,
		}
		global_vec_pool[i] = global_sync_pool
	}

}

type VECX struct {
	data      []float32
	dimension int
}

func newVEC2() interface{} {
	var vecx VECX
	vecx.dimension = 2
	vecx.data = make([]float32, 2, 2)
	return &vecx
}
func newVEC3() interface{} {
	var vecx VECX
	vecx.dimension = 3
	vecx.data = make([]float32, 3, 3)
	return &vecx
}
func newVEC4() interface{} {
	var vecx VECX
	vecx.dimension = 4
	vecx.data = make([]float32, 4, 4)
	return &vecx
}

// this is the only function that you should new-a-VECX
// the data inside the result is dirty, don't think they are zero-initialized
func GetVECX(di int) *VECX {
	return global_vec_pool[di].Get().(*VECX)
}

// when you dont need one VECX anymore, you should call this function
func DontNeedVECXAnyMore(vecx *VECX) {
	global_vec_pool[vecx.dimension].Put(vecx)
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
func (self *VECX) Add(other *VECX) *VECX {
	if !self.checkHomotype(other) {
		return nil
	}
	res := GetVECX(self.dimension)
	for i := 0; i != self.dimension; i++ {
		res.data[i] = self.data[i] + other.data[i]
	}
	return res
}

// vec math add in place, will store the result in self
// vec3 and vec4 cannot add, so return with nothing changed on that condition
func (self *VECX) Add_InPlace(other *VECX) {
	if !self.checkHomotype(other) {
		return
	}
	for i := 0; i != self.dimension; i++ {
		self.data[i] = self.data[i] + other.data[i]
	}
}

// vec math scale in place, will store the result in self
func (self *VECX) Scale_InPlace(factor float32) {
	for i := 0; i != self.dimension; i++ {
		self.data[i] = self.data[i] * factor
	}
}

// two vec interpolation
// t should be [0,1], but any other real number is supported
func (self *VECX) Interpolation(other *VECX, t float32) *VECX {
	if !self.checkHomotype(other) {
		return nil
	}
	res := GetVECX(self.dimension)
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
	temp_vecx_list := make([]*VECX, self.dimension, self.dimension)
	for i := 0; i != self.dimension; i++ {
		temp_vecx_list[i] = GetVECX(self.dimension)
		temp_vecx_list[i].GrabColToVec(other, i+1)
		temp_vecx_list[i].Scale_InPlace(self.GetIndexValue(i))
	}
	for i := 0; i != self.dimension; i++ {
		var temp_value float32
		for j := 0; j != self.dimension; j++ {
			temp_value += temp_vecx_list[j].GetIndexValue(i)
		}
		self.SetIndexValue(i, temp_value)
	}
	for i := 0; i != self.dimension; i++ {
		DontNeedVECXAnyMore(temp_vecx_list[i])
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
