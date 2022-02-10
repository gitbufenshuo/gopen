package matmath

import (
	"fmt"
	"math"
	"strings"

	"github.com/gitbufenshuo/gopen/help"
)

//
var VecZ = CreateVec4(0, 0, 1, 1)

//

type Vec4 struct {
	data [4]float32
}

func CreateVec4(x, y, z, w float32) Vec4 {
	return Vec4{
		data: [4]float32{x, y, z, w},
	}
}
func CreateVec4FromStr(data string) Vec4 {
	segs := strings.Split(data, ",")
	return Vec4{
		data: [4]float32{
			help.Str2Float32(segs[0]),
			help.Str2Float32(segs[1]),
			help.Str2Float32(segs[2]),
			help.Str2Float32(segs[3])},
	}
}
func (vec4 *Vec4) PrettyShow(reason string) {
	fmt.Printf("%s:(", reason)
	for i := 0; i != 4; i++ {
		fmt.Printf(" %v ", vec4.data[i])
	}
	fmt.Printf(")\n")
}
func (vec4 *Vec4) GetIndexValue(index int) float32 {
	return vec4.data[index]
}
func (vec4 *Vec4) SetIndexValue(index int, value float32) {
	vec4.data[index] = value
}
func (vec4 *Vec4) GetValue2() (float32, float32) {
	return vec4.data[0], vec4.data[1]
}
func (vec4 *Vec4) GetValue3() (float32, float32, float32) {
	return vec4.data[0], vec4.data[1], vec4.data[2]
}
func (vec4 *Vec4) GetValue4() (float32, float32, float32, float32) {
	return vec4.data[0], vec4.data[1], vec4.data[2], vec4.data[3]
}

func (vec4 *Vec4) SetX(value1 float32) {
	vec4.data[0] = value1
}

func (vec4 *Vec4) SetY(value1 float32) {
	vec4.data[1] = value1
}

func (vec4 *Vec4) SetZ(value1 float32) {
	vec4.data[2] = value1
}

func (vec4 *Vec4) SetValue2(value1, value2 float32) {
	vec4.data[0] = value1
	vec4.data[1] = value2
}

func (vec4 *Vec4) AddIndexValue(index int, value float32) {
	vec4.data[index] += value
}
func (vec4 *Vec4) Add2_InPlace(vec1, vec2 *Vec4) {
	for i := 0; i != 4; i++ {
		vec4.data[i] = vec1.data[i] + vec2.data[i]
	}
}

func (vec4 *Vec4) Sub(other *Vec4) Vec4 {
	var res Vec4
	for i := 0; i != 4; i++ {
		res.data[i] = vec4.data[i] - other.data[i]
	}
	return res
}

func (vec4 *Vec4) SetValue3(value1, value2, value3 float32) {
	vec4.data[0] = value1
	vec4.data[1] = value2
	vec4.data[2] = value3
}
func (vec4 *Vec4) SetValue4(value1, value2, value3, value4 float32) {
	vec4.data[0] = value1
	vec4.data[1] = value2
	vec4.data[2] = value3
	vec4.data[3] = value4
}
func (vec4 *Vec4) Clone(other *Vec4) {
	vec4.data[0] = other.data[0]
	vec4.data[1] = other.data[1]
	vec4.data[2] = other.data[2]
	vec4.data[3] = other.data[3]
}

// left X right
func Vec3Cross(left, right *Vec4) Vec4 {
	var res Vec4
	(res.data)[0] = (left.data)[1]*(right.data)[2] - (left.data)[2]*(right.data)[1]
	(res.data)[1] = (left.data)[2]*(right.data)[0] - (left.data)[0]*(right.data)[2]
	(res.data)[2] = (left.data)[0]*(right.data)[1] - (left.data)[1]*(right.data)[0]
	return res
}

// left product right 内积
func Vec3Pro(left, right *Vec4) float32 {
	return left.data[0]*right.data[0] + left.data[1]*right.data[1] + left.data[2]*right.data[2]
}

// vec3 模
func (vec3 *Vec4) Length() float32 {
	return help.Sqrt(vec3.data[0]*vec3.data[0] + vec3.data[1]*vec3.data[1] + vec3.data[2]*vec3.data[2])
}

// 向量夹角 cos
func Vec3CosTheta(left, right *Vec4) float32 {
	return Vec3Pro(left, right) / (left.Length() * right.Length())
}

// 归一化
func (vec4 *Vec4) Vec3Normalize() {
	weight := math.Sqrt(float64(vec4.data[0]*vec4.data[0] + vec4.data[1]*vec4.data[1] + vec4.data[2]*vec4.data[2]))
	if math.Abs(weight) < 0.000001 {
		return
	}
	for idx := range vec4.data {
		vec4.data[idx] = vec4.data[idx] / float32(weight)
	}
}

// grab one col from matx to store in self
func (vec4 *Vec4) GrabColToVec(mat4 *MAT4, col int) {
	for i := 0; i != 4; i++ {
		vec4.data[i] = mat4.GetEleByRowAndCol(i+1, col)
	}
}

// vec4 = mat4 * vec4
func (vec4 *Vec4) RightMul_InPlace(mat4 *MAT4) {
	var resarr [4]float32
	for rowMat := 0; rowMat != 4; rowMat++ {
		for vecidx := 0; vecidx != 4; vecidx++ {
			resarr[rowMat] += vec4.GetIndexValue(vecidx) * mat4.GetEleByRowAndCol(rowMat+1, vecidx+1)
		}
	}
	for idx := range resarr {
		vec4.data[idx] = resarr[idx]
	}
	return
}

// res = mat * vec4
func (vec4 Vec4) LeftMulMAT(mat MAT4) Vec4 {
	x, y, z, w := vec4.data[0], vec4.data[1], vec4.data[2], vec4.data[3]
	//
	return CreateVec4(
		x*mat.GetEleByRowAndCol(1, 1)+y*mat.GetEleByRowAndCol(1, 2)+z*mat.GetEleByRowAndCol(1, 3)+w*mat.GetEleByRowAndCol(1, 4),
		x*mat.GetEleByRowAndCol(2, 1)+y*mat.GetEleByRowAndCol(2, 2)+z*mat.GetEleByRowAndCol(2, 3)+w*mat.GetEleByRowAndCol(2, 4),
		x*mat.GetEleByRowAndCol(3, 1)+y*mat.GetEleByRowAndCol(3, 2)+z*mat.GetEleByRowAndCol(3, 3)+w*mat.GetEleByRowAndCol(3, 4),
		x*mat.GetEleByRowAndCol(4, 1)+y*mat.GetEleByRowAndCol(4, 2)+z*mat.GetEleByRowAndCol(4, 3)+w*mat.GetEleByRowAndCol(4, 4),
	)
}

func (vec4 *Vec4) InterpolationInplaceUnsafe(other *Vec4, t float32) {
	var res Vec4
	for i := 0; i != 4; i++ {
		res.data[i] = (1-t)*vec4.data[i] + t*other.data[i]
		vec4.data[i] = res.data[i]
	}
}

// generate mat4
// target : the target point coord
func LookAtFrom4(point, target, up *Vec4) MAT4 {
	var left MAT4
	var right MAT4
	left.ToIdentity()
	right.ToIdentity()

	// first lets calculate the camera-z and camera-x and camera-y
	// camera-z
	camera_z := point.Sub(target)
	camera_z.Vec3Normalize()
	// camera-x
	camera_x := Vec3Cross(up, &camera_z)
	camera_x.Vec3Normalize()

	// camera-y
	camera_y := Vec3Cross(&camera_z, &camera_x)
	camera_y.Vec3Normalize()
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
